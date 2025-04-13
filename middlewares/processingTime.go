package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Timing() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Microseconds()
		log.Printf("[%s] %s took %d µs", c.Method(), c.Path(), duration)

		return err
	}
}
