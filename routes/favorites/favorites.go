package favorites

import (
	"ChaiwalaBackend/db"
	"ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
)

func BuildRouter(app *fiber.App, dbConn *db.Queries) *fiber.Router {
	favoriteRouter := app.Group("/favorites")

	favoriteRouter.Post("/", buildFavoriteRecipe(dbConn))
	favoriteRouter.Delete("/", buildUnfavoriteRecipe(dbConn))
	favoriteRouter.Get("/", buildListUserFavorites(dbConn))

	return &favoriteRouter
}

func buildFavoriteRecipe(dbConn *db.Queries) fiber.Handler {
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

// buildUnfavoriteRecipe handles the logic for unfavoriting a recipe.
func buildUnfavoriteRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		unfavBody := new(Favorite)
		if err := c.Bind().JSON(unfavBody); err != nil {
			return err
		}

		// Unfavorite the recipe
		err := dbConn.UnfavoriteRecipe(c.Context(), db.UnfavoriteRecipeParams{UserID: unfavBody.RecipeID})
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

// buildListUserFavorites retrieves a list of all recipes favorited by the user.
func buildListUserFavorites(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userId").(int32) // Assuming user ID is available in Locals
		favorites, err := dbConn.ListUserFavorites(c.Context(), userID)
		if err != nil {
			c.Status(500)
			return c.JSON(routes.Error{
				Message:   "Could not retrieve favorites",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		return c.Status(200).JSON(favorites)
	}
}
