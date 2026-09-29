package main

import (
	"context"
	"log"

	"campus-lost-found-api/app/repository"
	"campus-lost-found-api/app/service"
	"campus-lost-found-api/config"
	"campus-lost-found-api/database"
	"campus-lost-found-api/helper"
	"campus-lost-found-api/middleware"
	"campus-lost-found-api/route"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := database.NewPool(
		ctx,
		cfg.DatabaseURL,
		cfg.DBMinConns,
		cfg.DBMaxConns,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	logger := config.NewLogger()
	defer logger.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: helper.ErrorHandler,
	})

	v := validator.New()

	psRaw := map[string][]string{}

	rows, err := pool.Query(
		ctx,
		`SELECT role_name, permission_name FROM role_permissions`,
	)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var role, perm string

		_ = rows.Scan(&role, &perm)

		psRaw[role] = append(psRaw[role], perm)
	}

	rows.Close()

	ps := helper.NewPermissionSet(psRaw)

	users := repository.NewUserRepository(pool)
	tokens := repository.NewTokenRepository(pool)
	itemsRepo := repository.NewItemRepository(pool)
	claimsRepo := repository.NewClaimRepository(pool)

	auth := service.NewAuthService(
		users,
		tokens,
		cfg,
		v,
	)

	items := service.NewItemService(
		itemsRepo,
		claimsRepo,
		v,
	)

	claims := service.NewClaimService(
		claimsRepo,
		itemsRepo,
		v,
	)

	admin := service.NewAdminService(users)

	app.Use(middleware.RequestID())
	app.Use(middleware.AccessLog(logger))
	app.Use(middleware.JSONOnly())

	route.Register(
		app,
		cfg,
		ps,
		auth,
		items,
		claims,
		admin,
	)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]any{
			"success": true,
			"message": "API is running",
		})
	})

	log.Printf("server running on :%s", cfg.Port)

	log.Fatal(app.Listen(":" + cfg.Port))
}