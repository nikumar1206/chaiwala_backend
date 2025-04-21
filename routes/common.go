package routes

import "github.com/gofiber/fiber/v3"

type Error struct {
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

func SendErrorResponse(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(
		Error{
			Message:   message,
			RequestId: c.GetRespHeader("X-Request-ID"),
		},
	)
}
