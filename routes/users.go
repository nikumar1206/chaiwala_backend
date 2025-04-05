package routes

import (
	"ChaiwalaBackend/db"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func BuildUsersRouter(app *fiber.App, db *db.Queries) *fiber.Router {
	userRouter := app.Group("/users")
	userRouter.Get("", buildGetUserHandler(db))
	return &userRouter
}

type Body struct {
	Slideshow struct {
		Author string `json:"author"`
		Date   string `json:"date"`
		Title  string `json:"title"`
	} `json:"slideshow"`
}

func buildGetUserHandler(db *db.Queries) fiber.Handler {

	return func(c fiber.Ctx) error {
		usr, err := db.GetUser(c.Context(), 1)
		if err != nil {
			fmt.Print(err)
			return c.JSON(Error{
				Message:   "Could not find requested user",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(usr)
	}

}
