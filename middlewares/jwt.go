// Middleware that validates JWT, skipping /login and /register
package middlewares

import (
	"net/http"
	"strings"

	jwtD "ChaiwalaBackend/jwt"

	"github.com/gofiber/fiber/v3"
)

func JWT() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()
		if path == "/users/login" || path == "/users/register" {
			return c.Next()
		}

		tokenStr := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		claims, err := jwtD.ValidateToken(tokenStr)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
