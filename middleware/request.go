package middleware

import (
	"campus-lost-found-api/config"
	"campus-lost-found-api/helper"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get("X-Request-Id")

		if id == "" {
			id = uuid.NewString()
		}

		c.Locals("request_id", id)
		c.Set("X-Request-Id", id)

		return c.Next()
	}
}

func JSONOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost ||
			c.Method() == fiber.MethodPut ||
			c.Method() == fiber.MethodPatch {

			ct := c.Get(fiber.HeaderContentType)

			if ct != "application/json" {
				return helper.Unsupported(
					"Content-Type harus application/json",
				)
			}
		}

		return c.Next()
	}
}

func AccessLog(l *config.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		requestID, _ := c.Locals("request_id").(string)

		l.Request(map[string]any{
			"request_id":  requestID,
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"duration_ms": time.Since(start).Milliseconds(),
		})

		return err
	}
}