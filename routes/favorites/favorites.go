package favorites

import (
	"net/http"
	"strconv"

	"ChaiwalaBackend/db"
	"ChaiwalaBackend/routes"
	common "ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
)

func BuildRouter(app *fiber.App, dbConn *db.Queries) *fiber.Router {
	favoriteRouter := app.Group("/favorites")

	favoriteRouter.Post("", favoriteRecipe(dbConn))
	favoriteRouter.Delete("/:favoriteId", unfavoriteRecipe(dbConn))

	return &favoriteRouter
}

func favoriteRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		favBody := new(Favorite)
		if err := c.Bind().JSON(favBody); err != nil {
			return err
		}

		// Favorite the recipe
		err := dbConn.FavoriteRecipe(c.Context(), db.FavoriteRecipeParams{UserID: favBody.UserID, RecipeID: favBody.RecipeID})
		if err != nil {
			c.Status(500)
			return c.JSON(routes.Error{
				Message:   "Could not favorite the recipe",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"Message": "Recipe favorited successfully",
		})
	}
}

// UnfavoriteRecipe handles the logic for unfavoriting a recipe.
func unfavoriteRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		favoriteID, err := strconv.Atoi(c.Params("favoriteId"))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid Favorite ID.")
		}

		// Unfavorite the recipe
		err = dbConn.UnfavoriteRecipe(c.Context(), db.UnfavoriteRecipeParams{UserID: int32(favoriteID)})
		if err != nil {
			c.Status(500)
			return c.JSON(routes.Error{
				Message:   "Could not unfavorite the recipe",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"Message": "Recipe unfavorited successfully",
		})
	}
}
