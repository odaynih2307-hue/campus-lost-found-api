package service

import (
	"campus-lost-found-api/app/model"
	"campus-lost-found-api/app/repository"
	"campus-lost-found-api/helper"
	"campus-lost-found-api/middleware"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

type ItemService struct {
	items  repository.ItemRepository
	claims repository.ClaimRepository
	v      *validator.Validate
}

func NewItemService(i repository.ItemRepository, c repository.ClaimRepository, v *validator.Validate) *ItemService {
	return &ItemService{i, c, v}
}
func (s *ItemService) List(c *fiber.Ctx) error {
	format, err := helper.Negotiate(c, "application/json", "text/csv")
	if err != nil {
		return err
	}
	q, err := parseCursorQuery(c)
	if err != nil {
		return err
	}
	u, _ := middleware.CurrentUser(c)
	if u.Role == "user" {
		q.OwnerID = u.UserID
	}
	rows, more, err := s.items.List(c.UserContext(), q)
	if err != nil {
		return helper.Internal(err)
	}
	if format == "text/csv" {
		return helper.WriteItemsCSV(c, rows)
	}
	meta := map[string]any{"limit": q.Limit, "has_more": more}
	if more && len(rows) > 0 {
		last := rows[len(rows)-1]
		meta["next_cursor"] = encodeCursor(last.CreatedAt, last.ID)
	}
	return helper.Success(c, 200, "daftar item berhasil diambil", rows, meta)
}
func parseCursorQuery(c *fiber.Ctx) (model.ListQuery, error) {
	limit := 10
	if v := c.QueryInt("limit", 10); v > 0 && v <= 50 {
		limit = v
	} else if c.Query("limit") != "" {
		return model.ListQuery{}, helper.BadRequest("limit harus 1 sampai 50")
	}
	q := model.ListQuery{Limit: limit, Search: strings.TrimSpace(c.Query("search")), Status: strings.TrimSpace(c.Query("status")), Category: strings.TrimSpace(c.Query("category"))}
	if cur := c.Query("cursor"); cur != "" {
		b, err := base64.RawURLEncoding.DecodeString(cur)
		if err != nil {
			return q, helper.BadRequest("cursor tidak valid")
		}
		parts := strings.Split(string(b), "|")
		if len(parts) != 2 {
			return q, helper.BadRequest("cursor tidak valid")
		}
		ns, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return q, helper.BadRequest("cursor tidak valid")
		}
		q.CursorCreated = timePtr(time.Unix(0, ns))
		q.CursorID, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return q, helper.BadRequest("cursor tidak valid")
		}
	}
	if q.Status != "" && q.Status != "lost" && q.Status != "found" && q.Status != "resolved" {
		return q, helper.BadRequest("status tidak valid")
	}
	return q, nil
}
func timePtr(t time.Time) *time.Time { return &t }
func encodeCursor(t time.Time, id int64) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d|%d", t.UnixNano(), id)))
}
func (s *ItemService) Get(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	i, err := s.items.FindByID(c.UserContext(), id)
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("item tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	u, _ := middleware.CurrentUser(c)
	if u.Role == "user" && i.ReportedBy != u.UserID {
		return helper.Forbidden("Anda tidak memiliki akses ke item ini")
	}
	return helper.Success(c, 200, "item berhasil diambil", i, nil)
}
func (s *ItemService) Create(c *fiber.Ctx) error {
	var req model.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	u, _ := middleware.CurrentUser(c)
	i, err := s.items.Create(c.UserContext(), model.Item{Title: req.Title, Description: req.Description, Category: req.Category, Location: req.Location, Status: req.Status, ReportedBy: u.UserID})
	if err != nil {
		return helper.Internal(err)
	}
	c.Set(fiber.HeaderLocation, fmt.Sprintf("/api/v1/items/%d", i.ID))
	return helper.Success(c, 201, "item berhasil dibuat", i, nil)
}
func (s *ItemService) Put(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	var req model.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	i, err := s.items.FindByID(c.UserContext(), id)
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("item tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	u, _ := middleware.CurrentUser(c)
	if !CanModifyItem(u.Role, i.ReportedBy, u.UserID) {
		return helper.Forbidden("Anda bukan pemilik item")
	}
	i.Title = req.Title
	i.Description = req.Description
	i.Category = req.Category
	i.Location = req.Location
	i.Status = req.Status
	i, err = s.items.Update(c.UserContext(), i)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "item berhasil diperbarui", i, nil)
}
func (s *ItemService) Patch(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	var req model.PatchItemRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	i, err := s.items.FindByID(c.UserContext(), id)
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("item tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	u, _ := middleware.CurrentUser(c)
	if !CanModifyItem(u.Role, i.ReportedBy, u.UserID) {
		return helper.Forbidden("Anda bukan pemilik item")
	}
	ApplyPatch(&i, req)
	i, err = s.items.Update(c.UserContext(), i)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "item berhasil diperbarui sebagian", i, nil)
}
func (s *ItemService) Delete(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	i, err := s.items.FindByID(c.UserContext(), id)
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("item tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	u, _ := middleware.CurrentUser(c)
	if !CanModifyItem(u.Role, i.ReportedBy, u.UserID) {
		return helper.Forbidden("Anda bukan pemilik item")
	}
	if err := s.items.Delete(c.UserContext(), id); err != nil {
		return helper.Internal(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
