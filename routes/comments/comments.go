package comments

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
	commentRouter := app.Group("/comments")

	commentRouter.Get("/:recipeId", listRecipeComments(dbConn))
	commentRouter.Get("/user/:userId", listUserComments(dbConn))
	commentRouter.Post("/", createComment(dbConn))
	commentRouter.Put("/:id", updateComment(dbConn))
	commentRouter.Delete("/:id", deleteComment(dbConn))

	return &commentRouter
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

func createComment(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		var comment CreateCommentBody
		if err := c.Bind().JSON(&comment); err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid input",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		createdComment, err := dbConn.AddComment(c.Context(), db.AddCommentParams{
			RecipeID: pgtype.Int4{Int32: comment.RecipeID, Valid: true},
			UserID:   pgtype.Int4{Int32: comment.UserID, Valid: true},
			Comment:  comment.Comment,
		})
		if err != nil {
			fmt.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to create comment",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.Status(http.StatusCreated).JSON(createdComment)
	}
}

func updateComment(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		commentID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid comment ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		var updateData UpdateCommentBody
		if err := c.Bind().JSON(&updateData); err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid input",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		err = dbConn.UpdateComment(c.Context(), db.UpdateCommentParams{
			ID:      int32(commentID),
			Comment: updateData.Comment,
		})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to update comment",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.SendStatus(http.StatusNoContent)
	}
}

func deleteComment(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		commentID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(common.Error{
				Message:   "Invalid comment ID",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}

		err = dbConn.DeleteComment(c.Context(), int32(commentID))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(common.Error{
				Message:   "Failed to delete comment",
				Context:   err.Error(),
				RequestId: c.GetRespHeader("X-Request-ID"),
			})
		}
		return c.SendStatus(http.StatusNoContent)
	}
}
