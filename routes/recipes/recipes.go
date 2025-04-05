package recipes

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
	recipeRouter := app.Group("/recipes")

	recipeRouter.Get("/", listPublicRecipes(dbConn))
	recipeRouter.Get("/:id", getRecipeByID(dbConn))
	recipeRouter.Post("/", createRecipe(dbConn))
	recipeRouter.Put("/:id", updateRecipe(dbConn))
	recipeRouter.Delete("/:id", deleteRecipe(dbConn))

	recipeRouter.Get("/user/:userId", listUserRecipes(dbConn))
	recipeRouter.Get("/favorites/user/:userId", listUserFavorites(dbConn))

	return &recipeRouter
}

func listPublicRecipes(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		recipes, err := dbConn.ListPublicRecipes(c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to fetch recipes",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(recipes)
	}
}

func getRecipeByID(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid recipe ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		recipe, err := dbConn.GetRecipe(c.Context(), int32(id))
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(common.Error{
				Message:   "Recipe not found",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(recipe)
	}
}

func createRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		var r CreateRecipeBody
		if err := c.Bind().JSON(&r); err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid input",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		recipe, err := dbConn.CreateRecipe(c.Context(), db.CreateRecipeParams{
			UserID:          pgtype.Int4{Int32: r.UserID, Valid: true},
			Title:           r.Title,
			Description:     r.Description,
			Instructions:    r.Instructions,
			ImageUrl:        r.ImageURL,
			PrepTimeMinutes: pgtype.Int4{Int32: r.PrepTimeMinutes, Valid: true},
			BrewTimeMinutes: pgtype.Int4{Int32: r.BrewTimeMinutes, Valid: true},
			Servings:        pgtype.Int4{Int32: r.Servings, Valid: true},
			IsPublic:        pgtype.Bool{Bool: r.IsPublic, Valid: true},
		})
		if err != nil {
			fmt.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to create recipe",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.Status(http.StatusCreated).JSON(recipe)
	}
}

func updateRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid recipe ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		var r UpdateRecipeBody
		if err := c.Bind().JSON(&r); err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid input",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		err = dbConn.UpdateRecipe(c.Context(), db.UpdateRecipeParams{
			ID:              int32(id),
			Title:           r.Title,
			Description:     r.Description,
			Instructions:    r.Instructions,
			ImageUrl:        r.ImageURL,
			PrepTimeMinutes: pgtype.Int4{Int32: r.PrepTimeMinutes, Valid: true},
			BrewTimeMinutes: pgtype.Int4{Int32: r.BrewTimeMinutes, Valid: true},
			Servings:        pgtype.Int4{Int32: r.Servings, Valid: true},
			IsPublic:        pgtype.Bool{Bool: r.IsPublic, Valid: true},
		})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to update recipe",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.SendStatus(http.StatusNoContent)
	}
}

func deleteRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid recipe ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		err = dbConn.DeleteRecipe(c.Context(), int32(id))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to delete recipe",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.SendStatus(http.StatusNoContent)
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
		userID, err := strconv.Atoi(c.Params("userId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid user ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		recipes, err := dbConn.ListUserFavorites(c.Context(), int32(userID))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to fetch user favorites",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(recipes)
	}
}
