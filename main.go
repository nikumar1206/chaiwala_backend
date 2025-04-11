package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"ChaiwalaBackend/db"
	"ChaiwalaBackend/middlewares"
	"ChaiwalaBackend/routes/comments"
	"ChaiwalaBackend/routes/favorites"
	"ChaiwalaBackend/routes/recipes"
	"ChaiwalaBackend/routes/users"

	_ "net/http/pprof"

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
	app.Use(requestid.New())
	app.Use(middlewares.JWT())

	conn, err := pgx.Connect(context.Background(), "user=nikhil dbname=chaiwala sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	dbConn := db.New(conn)

	users.BuildAuthRouter(app, dbConn)
	users.BuildRouter(app, dbConn)
	recipes.BuildRouter(app, dbConn)
	comments.BuildRouter(app, dbConn)
	favorites.BuildRouter(app, dbConn)

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
	app.Listen(ac.port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}
