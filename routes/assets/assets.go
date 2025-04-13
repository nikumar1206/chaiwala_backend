package assets

import (
	"fmt"
	"net/textproto"
	"net/url"

	"ChaiwalaBackend/clients/s3"

	"github.com/gofiber/fiber/v3"
)

func BuildRouter(app *fiber.App, s3Client s3.S3Client) *fiber.Router {
	assetsRouter := app.Group("/assets")

	uploadItem(s3Client)
	assetsRouter.Post("", uploadItem(s3Client))
	return &assetsRouter
}

func uploadItem(s3Client s3.S3Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		contentType := getContentType(file.Header)

		f, err := file.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		filename := url.PathEscape(file.Filename)
		s3Path := fmt.Sprintf("images/%s", filename)

		err = s3Client.Upload(
			c.Context(),
			s3Path,
			f,
			contentType,
		)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"Message": "Accepted File",
		})
	}
}

func getContentType(headers textproto.MIMEHeader) string {
	contentType := "application/octet-stream"
	if types := headers["Content-Type"]; len(types) > 0 {
		contentType = types[0]
	}
	return contentType
}
