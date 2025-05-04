package comments

import (
	"log/slog"
	"net/http"
	"strconv"

	"ChaiwalaBackend/db"
	common "ChaiwalaBackend/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

func BuildRouter(app *fiber.App, dbConn *db.Queries) *fiber.Router {
	commentRouter := app.Group("/comments")

	commentRouter.Post("", createComment(dbConn))
	commentRouter.Put("/:commentId", updateComment(dbConn))
	commentRouter.Delete("/:commentId", deleteComment(dbConn))

	return &commentRouter
}

func createComment(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		var comment CreateCommentBody
		if err := c.Bind().JSON(&comment); err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid Input")
		}

		createdComment, err := dbConn.AddComment(c.Context(), db.AddCommentParams{
			RecipeID: pgtype.Int4{Int32: comment.RecipeID, Valid: true},
			UserID:   pgtype.Int4{Int32: comment.UserID, Valid: true},
			Comment:  comment.Comment,
		})
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return common.SendErrorResponse(c, http.StatusBadRequest, "Failed to create comment.")
		}
		return c.Status(http.StatusCreated).JSON(createdComment)
	}
}

func updateComment(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		commentID, err := strconv.Atoi(c.Params("commentId"))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid Comment Id.")
		}

		var updateData UpdateCommentBody
		if err := c.Bind().JSON(&updateData); err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid input.")
		}

		err = dbConn.UpdateComment(c.Context(), db.UpdateCommentParams{
			ID:      int32(commentID),
			Comment: updateData.Comment,
		})
		if err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Failed to update comment")
		}
		return c.SendStatus(http.StatusNoContent)
	}
}

func deleteComment(dbConn *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		commentID, err := strconv.Atoi(c.Params("commentId"))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusBadRequest, "Invalid Comment Id.")
		}

		err = dbConn.DeleteComment(c.Context(), int32(commentID))
		if err != nil {
			return common.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete comment")
		}
		return c.SendStatus(http.StatusNoContent)
	}
}
