package assets

import (
	"fmt"
	"log/slog"
	"net/textproto"
	"time"

	"ChaiwalaBackend/clients/s3"
	"ChaiwalaBackend/routes"
	"ChaiwalaBackend/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var DEFAULT_CONTENT_TYPE string = "application/octet-stream"

func BuildRouter(app *fiber.App, s3Client s3.S3Client) *fiber.Router {
	fileRouter := app.Group("/files")

	fileRouter.Post("", uploadItem(s3Client))
	fileRouter.Get("/:fileId", getItem(s3Client))
	return &fileRouter
}

func uploadItem(s3Client s3.S3Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return routes.SendErrorResponse(c, fiber.StatusUnprocessableEntity, "Request must include a Mulipart file under attribute 'file'")
		}

		contentType := getContentType(file.Header)

		f, err := file.Open()
		if err != nil {
			return routes.SendErrorResponse(c, fiber.StatusUnprocessableEntity, "Unable to open the provided file. Please make sure its complete.")
		}
		defer utils.LogThrowable(c.Context(), f.Close())

		fileId := uuid.NewString()
		s3Path := "images/" + fileId

		err = s3Client.Upload(
			c.Context(),
			s3Path,
			f,
			contentType,
		)
		if err != nil {
			return routes.SendErrorResponse(c, 500, "Unable to upload the file. Please try again later.")
		}

		return c.JSON(fiber.Map{
			"Message": "Accepted File",
			"FileId":  fileId,
		})
	}
}

func getItem(s3Client s3.S3Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		fileId := c.Params("fileId")

		if fileId == "" {
			return routes.SendErrorResponse(c, fiber.StatusUnprocessableEntity, "No fileId provided")
		}
		slog.InfoContext(c.Context(), "request for file", slog.String("fileId", fileId))
		key := fmt.Sprintf("images/%s", fileId)
		startTime := time.Now()
		resp, err := s3Client.Download(c.Context(), key)
		if err != nil {
			slog.ErrorContext(c.Context(), err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to download file from S3")
		}

		slog.InfoContext(c.Context(), "time since download", slog.String("duration", time.Since(startTime).String()))

		if resp.ContentType != nil {
			c.Set("Content-Type", *resp.ContentType)
		}
		slog.InfoContext(c.Context(), "resp.content=", slog.String("contentType", *resp.ContentType))
		if resp.ContentLength != nil && *resp.ContentLength > 0 {
			c.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
		}
		if resp.ContentDisposition != nil {
			c.Set("Content-Disposition", *resp.ContentDisposition)
		}
		c.Response().SetBodyStream(resp.Body, -1)

		return nil
	}
}

func getContentType(headers textproto.MIMEHeader) string {
	if types := headers["Content-Type"]; len(types) > 0 {
		return types[0]
	}
	return DEFAULT_CONTENT_TYPE
}
