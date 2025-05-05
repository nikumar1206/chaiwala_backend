package recipes

import (
	"log/slog"
	"net/http"
	"strconv"

	"ChaiwalaBackend/db"
	logger "ChaiwalaBackend/logging"
	common "ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func BuildRouter(app *fiber.App, conn *pgx.Conn, dbConn *db.Queries) *fiber.Router {
	recipeRouter := app.Group("/recipes")

	recipeRouter.Get("", listPublicRecipes(dbConn))
	recipeRouter.Get("/:recipeId", getRecipeByID(dbConn))
	recipeRouter.Post("", createRecipe(conn))
	recipeRouter.Put("/:recipeId", updateRecipe(conn))
	recipeRouter.Delete("/:recipeId", deleteRecipe(dbConn))

	recipeRouter.Get("/:recipeId/comments", listRecipeComments(dbConn))

	return &recipeRouter
}

func listPublicRecipes(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		qParams := c.Queries()

		offset := qParams["offset"]
		limit := qParams["limit"]

		if offset == "" {
			offset = "0"
		}

		if limit == "" {
			limit = "10"
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusUnprocessableEntity, "Invalid offset")
		}

		if offsetInt < 0 {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Offset must be greater than or equal to 0")
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusUnprocessableEntity, "Invalid limit")
		}
		if limitInt < 1 {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Limit must be greater than 0")
		}
		if limitInt > 500 {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Limit must be less than 500")
		}

		recipes, err := dbConn.ListPublicRecipesPaginated(c.Context(), db.ListPublicRecipesPaginatedParams{
			Limit:  int32(limitInt),
			Offset: int32(offsetInt),
		})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
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
		slog.InfoContext(c.Context(), "get recipe by id", slog.Int("id", id))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusUnprocessableEntity, "Invalid Request ID")
		}

		recipe, err := dbConn.GetRecipe(c.Context(), int32(id))
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusNotFound, "Recipe not found")
		}

		steps, err := dbConn.GetRecipeStepsByRecipe(c.Context(), pgtype.Int4{Int32: recipe.ID, Valid: true})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch recipe steps")
		}

		if steps == nil {
			steps = []db.RecipeStep{}
		}

		user, err := dbConn.GetUser(c.Context(), recipe.UserID.Int32)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user")
		}
		r := GetRecipe{
			Recipe:         recipe,
			Steps:          steps,
			ID:             recipe.ID,
			CommentsCount:  0,
			FavoritesCount: 0,
			CreatedBy:      user,
		}

		return c.JSON(r)
	}
}

func createRecipe(conn *pgx.Conn) fiber.Handler {
	return func(c fiber.Ctx) error {
		var r CreateRecipeBody
		if err := c.Bind().JSON(&r); err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid input",
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		tx, err := conn.Begin(c.Context())
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Sorry, something went wrong. Please try again later.")
		}

		defer func() {
			if err != nil {
				if rbErr := tx.Rollback(c.Context()); rbErr != nil {
					slog.ErrorContext(c.Context(), err.Error())
				}
			}
		}()

		q := db.New(tx)
		userId := c.Locals(logger.UserId).(int32)
		recipe, err := q.CreateRecipe(c.Context(), db.CreateRecipeParams{
			UserID:          pgtype.Int4{Int32: userId, Valid: true},
			Title:           r.Title,
			Description:     r.Description,
			Type:            int32(r.TeaType),
			AssetID:         r.AssetId,
			PrepTimeMinutes: pgtype.Int4{Int32: r.PrepTimeMinutes, Valid: true},
			Servings:        pgtype.Int4{Int32: r.Servings, Valid: true},
			IsPublic:        pgtype.Bool{Bool: r.IsPublic, Valid: true},
		})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create recipe")
		}
		slog.InfoContext(c.Context(), "scheduled recipe", slog.Int("recipeId", int(recipe.ID)))
		for _, step := range r.Steps {
			_, err := q.AddRecipeStep(c.Context(), db.AddRecipeStepParams{
				RecipeID:    pgtype.Int4{Int32: recipe.ID, Valid: true},
				StepNumber:  int32(step.StepNumber),
				Description: step.Description,
				AssetID:     pgtype.Text{String: step.AssetId, Valid: true},
			})
			if err != nil {
				slog.ErrorContext(c.Context(), err.Error())
				return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create recipe")
			}
		}
		slog.InfoContext(c.Context(), "scheduled steps")

		err = tx.Commit(c.Context())
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		}

		slog.InfoContext(c.Context(), "success")
		return c.Status(http.StatusCreated).JSON(recipe)
	}
}

func updateRecipe(conn *pgx.Conn) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("recipeId"))
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid recipe ID")
		}
		var r UpdateRecipeBody
		if err := c.Bind().JSON(&r); err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		}
		tx, err := conn.Begin(c.Context())
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Sorry, something went wrong. Please try again later.")
		}
		defer func() {
			if err != nil {
				slog.ErrorContext(c.Context(), err.Error())
				tx.Rollback(c.Context())
			}
		}()

		q := db.New(tx)

		err = q.UpdateRecipe(c.Context(), db.UpdateRecipeParams{
			ID:              int32(id),
			Title:           r.Title,
			Description:     r.Description,
			Type:            int32(r.TeaType),
			AssetID:         r.AssetID,
			PrepTimeMinutes: pgtype.Int4{Int32: r.PrepTimeMinutes, Valid: true},
			Servings:        pgtype.Int4{Int32: r.Servings, Valid: true},
			IsPublic:        pgtype.Bool{Bool: r.IsPublic, Valid: true},
		})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update recipe")
		}

		for _, step := range r.Steps {
			err := q.UpdateRecipeStep(c.Context(), db.UpdateRecipeStepParams{
				ID:          int32(step.ID),
				StepNumber:  int32(step.StepNumber),
				Description: step.Description,
				AssetID:     step.AssetID,
			})
			if err != nil {
				return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update recipe step")
			}
		}
		tx.Commit(c.Context())
		slog.InfoContext(c.Context(), "Recipe updated successfully")
		return c.SendStatus(http.StatusNoContent)
	}
}

func deleteRecipe(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("recipeId"))
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid recipe ID")
		}
		err = dbConn.DeleteRecipe(c.Context(), int32(id))
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete recipe")
		}
		slog.InfoContext(c.Context(), "Recipe deleted successfully")
		return c.SendStatus(http.StatusNoContent)
	}
}

func listRecipeComments(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		recipeID, err := strconv.Atoi(c.Params("recipeId"))
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid recipe ID")
		}

		comments, err := dbConn.ListComments(c.Context(), pgtype.Int4{Int32: int32(recipeID), Valid: true})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch comments")
		}
		return c.JSON(comments)
	}
}
