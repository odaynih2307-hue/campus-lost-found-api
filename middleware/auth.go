package middleware

import (
	"campus-lost-found-api/app/model"
	"campus-lost-found-api/config"
	"campus-lost-found-api/helper"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func RequireAuth(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get(fiber.HeaderAuthorization)
		parts := strings.Fields(h)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return helper.Unauthorized("token tidak ditemukan")
		}
		u, err := helper.ParseAccessToken(parts[1], cfg.JWTSecret, cfg.JWTIssuer)
		if err != nil {
			return err
		}
		c.Locals("auth_user", u)
		return c.Next()
	}
}
func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	u, ok := c.Locals("auth_user").(model.AuthUser)
	return u, ok
}
func RequirePermission(ps *helper.PermissionSet, p string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u, ok := CurrentUser(c)
		if !ok {
			return helper.Unauthorized("belum terautentikasi")
		}
		if !ps.Can(u.Role, p) {
			return helper.Forbidden("akses ditolak")
		}
		return c.Next()
	}
}
