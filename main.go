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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jackc/pgx/v5"
)

type AppConfig struct {
	PORT                  string
	AWS_REGION            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	S3_BUCKET_NAME        string
}

func newAppConfig() *AppConfig {
	return &AppConfig{
		PORT:                  ":" + os.Getenv("PORT"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3_BUCKET_NAME:        os.Getenv("S3_BUCKET_NAME"),
	}
}

func main() {
	ac := newAppConfig()

	// s3Client := s3Helper.New(context.Background(), ac.AWS_REGION, ac.S3_BUCKET_NAME)
	// err := s3Client.Upload(
	// 	context.Background(),
	// 	"test/test.jpeg",
	// 	readFile("test.jpeg"),
	// 	"image/jpeg",
	// )
	// if err != nil {
	// 	panic(err)
	// }

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(middlewares.Timing())
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
	app.Listen(ac.PORT, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}

func getS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return s3.NewFromConfig(cfg)
}

func readFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return content
}
