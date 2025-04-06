package users

import (
	"fmt"
	"net/http"
	"strconv"

	"ChaiwalaBackend/db"
	common "ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

func BuildRouter(app *fiber.App, dbConn *db.Queries) *fiber.Router {
	userRouter := app.Group("/users")

	userRouter.Get("/:userId", getUser(dbConn))
	userRouter.Get("/:userId/recipes", listUserRecipes(dbConn))
	userRouter.Get("/:userId/favorites", listUserFavorites(dbConn))
	userRouter.Get("/:userId/comments", listUserComments(dbConn))

	return &userRouter
}

func getUser(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		usr, err := dbConn.GetUser(c.Context(), 1)
		if err != nil {
			fmt.Print(err)
			return c.JSON(common.Error{
				Message:   "User not found.",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(usr)
	}
}

func listUserRecipes(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, err := strconv.Atoi(c.Params("userId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid user ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		recipes, err := dbConn.ListUserRecipes(c.Context(), pgtype.Int4{Int32: int32(userID), Valid: true})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to fetch user recipes",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(recipes)
	}
}

func listUserFavorites(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userId").(int32) // Assuming user ID is available in Locals
		favorites, err := dbConn.ListUserFavorites(c.Context(), userID)
		if err != nil {
			c.Status(500)
			return c.JSON(common.Error{
				Message:   "Could not retrieve favorites",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		return c.Status(200).JSON(favorites)
	}
}

func listUserComments(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, err := strconv.Atoi(c.Params("userId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid user ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		comments, err := dbConn.ListCommentsByUser(c.Context(), pgtype.Int4{Int32: int32(userID), Valid: true})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to fetch user comments",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(comments)
	}
}
