package middlewares

import (
	"context"
	"log/slog"

	logger "ChaiwalaBackend/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func SetContext() fiber.Handler {
	return func(c fiber.Ctx) error {
		slog.Debug("adding additional request context")
		requestId := uuid.NewString()
		ctx := c.Context()

		ctx = context.WithValue(ctx, logger.RequestId, requestId)
		ctx = context.WithValue(ctx, logger.Method, c.Method())
		ctx = context.WithValue(ctx, logger.Path, c.Path())
		ctx = context.WithValue(ctx, logger.SourceIP, c.IP())
		// set user related context settings in jwt middleware
		c.SetContext(ctx)
		c.Response().Header.Set("X-Request-ID", requestId)

		return c.Next()
	}
}
