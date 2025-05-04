package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	jwtD "ChaiwalaBackend/clients/jwt"
	logger "ChaiwalaBackend/logging"
	"ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
)

func JWT() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()

		if path == "/auth/login" || path == "/auth/register" {
			slog.InfoContext(c.Context(), "skipping on auth routes")
			return c.Next()
		}

		tokenStr := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		claims, err := jwtD.ValidateToken(c, tokenStr)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return routes.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		// set necessary contextvars
		c.Locals(logger.Email, claims.Email)
		c.Locals(logger.UserId, claims.UserID)
		c.Locals("token", tokenStr)
		c.Locals("claims", claims)

		ctx := c.Context()
		ctx = context.WithValue(ctx, logger.Email, claims.Email)
		ctx = context.WithValue(ctx, logger.UserId, claims.UserID)

		c.SetContext(ctx)

		return c.Next()
	}
}
