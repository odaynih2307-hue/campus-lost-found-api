package route

import (
	"campus-lost-found-api/app/service"
	"campus-lost-found-api/config"
	"campus-lost-found-api/helper"
	"campus-lost-found-api/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, cfg config.Config, ps *helper.PermissionSet, auth *service.AuthService, items *service.ItemService, claims *service.ClaimService, admin *service.AdminService) {
	api := app.Group("/api/v1")
	authr := api.Group("/auth")
	authr.Post("/register", auth.Register)
	authr.Post("/login", auth.Login)
	authr.Post("/refresh", auth.Refresh)
	authr.Post("/logout", auth.Logout)
	authr.Get("/me", middleware.RequireAuth(cfg), auth.Me)
	secure := api.Group("", middleware.RequireAuth(cfg))
	secure.Get("/items", middleware.RequirePermission(ps, "item:list"), items.List)
	secure.Get("/items/:id", middleware.RequirePermission(ps, "item:list"), items.Get)
	secure.Post("/items", middleware.RequirePermission(ps, "item:create"), items.Create)
	secure.Put("/items/:id", middleware.RequirePermission(ps, "item:list"), items.Put)
	secure.Patch("/items/:id", middleware.RequirePermission(ps, "item:list"), items.Patch)
	secure.Delete("/items/:id", middleware.RequirePermission(ps, "item:delete"), items.Delete)
	secure.Post("/items/:id/claims", middleware.RequirePermission(ps, "claim:create"), claims.Create)
	secure.Get("/items/:id/claims", middleware.RequirePermission(ps, "claim:list"), claims.List)
	secure.Patch("/claims/:id/status", middleware.RequirePermission(ps, "claim:manage"), claims.UpdateStatus)
	secure.Patch("/users/:id/role", middleware.RequirePermission(ps, "role:assign"), admin.AssignRole)
}
