package service

import (
	"campus-lost-found-api/app/model"
	"campus-lost-found-api/app/repository"
	"campus-lost-found-api/config"
	"campus-lost-found-api/helper"
	"campus-lost-found-api/middleware"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type AuthService struct {
	users  repository.UserRepository
	tokens repository.TokenRepository
	cfg    config.Config
	v      *validator.Validate
}

func NewAuthService(u repository.UserRepository, t repository.TokenRepository, c config.Config, v *validator.Validate) *AuthService {
	return &AuthService{u, t, c, v}
}
func (s *AuthService) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Internal(err)
	}
	u, err := s.users.Create(c.UserContext(), model.User{Username: req.Username, Email: req.Email, Password: hash})
	if errors.Is(err, repository.ErrDuplicate) {
		return helper.Conflict("username atau email sudah digunakan")
	}
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 201, "registrasi berhasil", u, nil)
}
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	u, err := s.users.FindByUsername(c.UserContext(), strings.TrimSpace(req.Username))
	if err != nil || !u.IsActive || !helper.CheckPassword(u.Password, req.Password) {
		return helper.Unauthorized("username atau password salah")
	}
	pair, err := s.issue(c.UserContext(), u)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "login berhasil", pair, nil)
}
func (s *AuthService) issue(ctx context.Context, u model.User) (model.TokenPair, error) {
	access, err := helper.IssueAccessToken(u, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.AccessTokenMinutes)
	if err != nil {
		return model.TokenPair{}, err
	}
	raw, err := helper.RandomToken()
	if err != nil {
		return model.TokenPair{}, err
	}
	exp := time.Now().Add(time.Duration(s.cfg.RefreshTokenDays) * 24 * time.Hour)
	if err := s.tokens.Create(ctx, u.ID, helper.SHA256Hex(raw), exp); err != nil {
		return model.TokenPair{}, err
	}
	return model.TokenPair{AccessToken: access, RefreshToken: raw, TokenType: "Bearer", ExpiresIn: s.cfg.AccessTokenMinutes * 60}, nil
}
func (s *AuthService) Refresh(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if err := helper.Validate(s.v, req); err != nil {
		return err
	}
	hash := helper.SHA256Hex(req.RefreshToken)
	st, err := s.tokens.FindActive(c.UserContext(), hash)
	if err != nil {
		return helper.Unauthorized("refresh token tidak valid atau sudah kedaluwarsa")
	}
	u, err := s.users.FindByID(c.UserContext(), st.UserID)
	if err != nil || !u.IsActive {
		return helper.Unauthorized("akun tidak dapat dipakai")
	}
	if err := s.tokens.Revoke(c.UserContext(), hash); err != nil {
		return helper.Internal(err)
	}
	pair, err := s.issue(c.UserContext(), u)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, 200, "token berhasil diperbarui", pair, nil)
}
func (s *AuthService) Logout(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err == nil && strings.TrimSpace(req.RefreshToken) != "" {
		_ = s.tokens.Revoke(c.UserContext(), helper.SHA256Hex(req.RefreshToken))
	}
	return helper.Success(c, 200, "logout berhasil", nil, nil)
}
func (s *AuthService) Me(c *fiber.Ctx) error {
	u, ok := middleware.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}
	user, err := s.users.FindByID(c.UserContext(), u.UserID)
	if err != nil {
		return helper.Unauthorized("user tidak ditemukan")
	}
	return helper.Success(c, 200, "profil berhasil diambil", user, nil)
}
