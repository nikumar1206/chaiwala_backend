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

	recipeRouter.Get("", listPublicRecipes(dbConn))
	recipeRouter.Get("/:recipeId", getRecipeByID(dbConn))
	recipeRouter.Post("", createRecipe(dbConn))
	recipeRouter.Put("/:recipeId", updateRecipe(dbConn))
	recipeRouter.Delete("/:recipeId", deleteRecipe(dbConn))

	recipeRouter.Get("/:recipeId/comments", listRecipeComments(dbConn))

	return &recipeRouter
}

func listPublicRecipes(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		recipes, err := dbConn.ListPublicRecipes(c.Context())
		if err != nil {
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch recipes")
		}
		if recipes == nil {
			return c.JSON([]db.Recipe{})
		}
		return c.JSON(recipes)
	}
}

func getRecipeByID(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("recipeId"))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusUnprocessableEntity, "Invalid Request ID")
		}

		recipe, err := dbConn.GetRecipe(c.Context(), int32(id))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusNotFound, "Recipe not found")
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
			ImageUrl:        r.ImageData,
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
		id, err := strconv.Atoi(c.Params("recipeId"))
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
		id, err := strconv.Atoi(c.Params("recipeId"))
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

func listRecipeComments(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		recipeID, err := strconv.Atoi(c.Params("recipeId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid recipe ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		comments, err := dbConn.ListComments(c.Context(), pgtype.Int4{Int32: int32(recipeID), Valid: true})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to fetch comments",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.JSON(comments)
	}
}
