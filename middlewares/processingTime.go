package middlewares

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Timing() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()
		slog.InfoContext(
			c.Context(),
			"handled request",
			slog.String("duration", time.Since(start).String()),
			slog.Int("statusCode", c.Response().StatusCode()),
		)

		return err
	}
}
