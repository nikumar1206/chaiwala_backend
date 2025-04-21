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
		claims, err := jwtD.ValidateToken(tokenStr)
		if err != nil {
			return routes.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		}
		c.Locals("username", claims.Username)
		c.Locals("userId", claims.UserID)

		ctx := c.Context()
		ctx = context.WithValue(ctx, logger.Username, claims.Username)
		ctx = context.WithValue(ctx, logger.UserId, claims.UserID)

		c.SetContext(ctx)

		return c.Next()
	}
}
