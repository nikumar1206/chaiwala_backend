package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"ChaiwalaBackend/clients/s3"
	"ChaiwalaBackend/db"
	logger "ChaiwalaBackend/logging"
	"ChaiwalaBackend/middlewares"
	"ChaiwalaBackend/routes/assets"
	"ChaiwalaBackend/routes/comments"
	"ChaiwalaBackend/routes/favorites"
	"ChaiwalaBackend/routes/recipes"
	"ChaiwalaBackend/routes/users"
	"ChaiwalaBackend/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type AppConfig struct {
	APP_ENV               string
	PORT                  string
	AWS_REGION            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	S3_BUCKET_NAME        string
	LOG_LEVEL             slog.Level
}

func newAppConfig() *AppConfig {
	logLevel := utils.Must(strconv.Atoi((os.Getenv("LOG_LEVEL"))))

	return &AppConfig{
		APP_ENV:               os.Getenv("APP_ENV"),
		PORT:                  ":" + os.Getenv("PORT"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3_BUCKET_NAME:        os.Getenv("S3_BUCKET_NAME"),
		LOG_LEVEL:             slog.Level(logLevel),
	}
}

func main() {
	ac := newAppConfig()

	logger := slog.New(logger.CustomHandler{Handler: getLoggerHandler(ac)})
	slog.SetDefault(logger)

	app := fiber.New()

	app.Use(middlewares.SetContext())
	app.Use(middlewares.Timing())
	app.Use(middlewares.JWT())

	conn := utils.Must(pgx.Connect(context.Background(), "user=nikhil dbname=chaiwala sslmode=verify-full"))
	//nolint:errcheck
	defer conn.Close(context.Background())

	dbConn := db.New(conn)
	s3Client := s3.New(context.Background(), ac.AWS_REGION, ac.S3_BUCKET_NAME)

	users.BuildAuthRouter(app, dbConn)
	users.BuildRouter(app, dbConn)
	recipes.BuildRouter(app, dbConn)
	comments.BuildRouter(app, dbConn)
	favorites.BuildRouter(app, dbConn)
	assets.BuildRouter(app, s3Client)

	app.Get("", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	routes := app.GetRoutes(true)

	for _, route := range routes {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}

	utils.LogThrowable(
		context.Background(),
		app.Listen(ac.PORT, fiber.ListenConfig{DisableStartupMessage: true}))
}

func getLoggerHandler(ac *AppConfig) slog.Handler {
	if ac.APP_ENV == "PRODUCTION" {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     ac.LOG_LEVEL,
			AddSource: true,
		})
	}

	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     ac.LOG_LEVEL,
		AddSource: true,
	})
}
