package service

import (
	"campus-lost-found-api/app/repository"
	"campus-lost-found-api/helper"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type AdminService struct{ users repository.UserRepository }

func NewAdminService(u repository.UserRepository) *AdminService { return &AdminService{u} }
func (s *AdminService) AssignRole(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)
	if err != nil {
		return err
	}
	var body struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if body.Role != "admin" && body.Role != "staff" && body.Role != "user" {
		return helper.BadRequest("role tidak valid")
	}
	if err := s.users.UpdateRole(c.UserContext(), id, body.Role); errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("user tidak ditemukan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "role berhasil diperbarui", map[string]any{"user_id": id, "role": body.Role}, nil)
}
