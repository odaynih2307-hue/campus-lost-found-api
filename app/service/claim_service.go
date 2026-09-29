package service

import (
	"campus-lost-found-api/app/model"
	"campus-lost-found-api/app/repository"
	"campus-lost-found-api/helper"
	"campus-lost-found-api/middleware"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ClaimService struct {
	claims repository.ClaimRepository
	items  repository.ItemRepository
	v      *validator.Validate
}

func NewClaimService(c repository.ClaimRepository, i repository.ItemRepository, v *validator.Validate) *ClaimService {
	return &ClaimService{c, i, v}
}
func (s *ClaimService) Create(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	if _, err := s.items.FindByID(c.UserContext(), id); errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("item tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	var req model.CreateClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	u, _ := middleware.CurrentUser(c)
	cl, err := s.claims.Create(c.UserContext(), model.Claim{ItemID: id, ClaimantID: u.UserID, Note: req.Note})
	if err != nil {
		return helper.Conflict("Anda sudah mengajukan klaim untuk item ini")
	}
	return helper.Success(c, 201, "klaim berhasil dibuat", cl, nil)
}
func (s *ClaimService) List(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	item, err := s.items.FindByID(c.UserContext(), id)
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("item tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	u, _ := middleware.CurrentUser(c)
	if u.Role == "user" && item.ReportedBy != u.UserID {
		return helper.Forbidden("Anda bukan pemilik item")
	}
	rows, err := s.claims.ListByItem(c.UserContext(), id)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "daftar klaim berhasil diambil", rows, nil)
}
func (s *ClaimService) UpdateStatus(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	var req model.UpdateClaimStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	cl, err := s.claims.UpdateStatus(c.UserContext(), id, req.Status)
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("klaim tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "status klaim berhasil diperbarui", cl, nil)
}
