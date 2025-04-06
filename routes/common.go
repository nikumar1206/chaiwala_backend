package routes

import "github.com/gofiber/fiber/v3"

type Error struct {
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	Context   string `json:"context,omitempty"` // todo(nick): remove once all error routes are using SendErrorResponse
}

func SendErrorResponse(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(
		Error{
			Message:   message,
			RequestId: c.GetRespHeader("X-Request-ID"),
		},
	)
}
