package favorites

import (
	"net/http"
	"strconv"

	"ChaiwalaBackend/db"
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

		err := dbConn.FavoriteRecipe(c.Context(), db.FavoriteRecipeParams{UserID: favBody.UserID, RecipeID: favBody.RecipeID})
		if err != nil {
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Could not favorite the recipe")
		}

		return c.Status(200).JSON(fiber.Map{
			"Message": "Recipe favorited successfully",
		})
	}
}

func unfavoriteRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		favoriteID, err := strconv.Atoi(c.Params("favoriteId"))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid Favorite ID.")
		}

		err = dbConn.UnfavoriteRecipe(c.Context(), db.UnfavoriteRecipeParams{UserID: int32(favoriteID)})
		if err != nil {
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Could not unfavorite the recipe")
		}

		return c.Status(200).JSON(fiber.Map{
			"Message": "Recipe unfavorited successfully",
		})
	}
}
