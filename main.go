package main

import (
	"ChaiwalaBackend/db"
	"ChaiwalaBackend/routes"
	"context"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jackc/pgx/v5"
)

type AppConfig struct {
	port string
}

func newAppConfig() *AppConfig {
	return &AppConfig{
		port: ":" + os.Getenv("PORT"),
	}
}

func main() {

	ac := newAppConfig()

	app := fiber.New()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=nikhil dbname=chaiwala sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	dbConn := db.New(conn)

	app.Use(requestid.New())

	routes.BuildUsersRouter(app, dbConn)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(ac.port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})

}
