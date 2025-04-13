package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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

	_ "net/http/pprof"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type AppConfig struct {
	PORT                  string
	AWS_REGION            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	S3_BUCKET_NAME        string
	LOG_LEVEL             slog.Level
}

func newAppConfig() *AppConfig {
	logLevel := Must(strconv.Atoi((os.Getenv("LOG_LEVEL"))))

	return &AppConfig{
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

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     ac.LOG_LEVEL,
		AddSource: true,
	})

	logger := slog.New(logger.FiberHandler{Handler: handler})
	slog.SetDefault(logger)

	app := fiber.New()

	app.Use(middlewares.SetContext())
	app.Use(middlewares.Timing())
	app.Use(middlewares.JWT())

	conn := Must(pgx.Connect(context.Background(), "user=nikhil dbname=chaiwala sslmode=verify-full"))
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
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	routes := app.GetRoutes(true)

	for _, route := range routes {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
	app.Listen(ac.PORT, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}
