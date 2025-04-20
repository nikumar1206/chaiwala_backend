package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwtD "ChaiwalaBackend/jwt"
	logger "ChaiwalaBackend/logging"

	"github.com/gofiber/fiber/v3"
)

func JWT() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()
		fmt.Println("called jwt")
		if path == "/auth/login" || path == "/auth/register" {
			return c.Next()
		}

		tokenStr := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		claims, err := jwtD.ValidateToken(tokenStr)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
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
