package middlewares

import (
	"context"

	logger "ChaiwalaBackend/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func SetContext() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestId := uuid.NewString()
		ctx := c.Context()

		ctx = context.WithValue(ctx, logger.RequestId, requestId)
		ctx = context.WithValue(ctx, logger.Method, c.Method())
		ctx = context.WithValue(ctx, logger.Path, c.Path())
		ctx = context.WithValue(ctx, logger.SourceIP, c.IP())

		c.SetContext(ctx)
		return c.Next()
	}
}
