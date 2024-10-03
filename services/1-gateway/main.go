package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error reading environment data %v", err)
	}

	port := os.Getenv("PORT")
	app := fiber.New(fiber.Config{
		BodyLimit:     5 * 1024 * 1024,
		CaseSensitive: true,
		StrictRouting: true,
		// Prefork:       true,
	})
	db_user := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", db_user, db_password, db_name, db_port)

	fmt.Println("Connected to postgres DB", dsn)

	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		Storage: postgres.New(postgres.Config{
			ConnectionURI: dsn,
		}),
	}))
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"),
		AllowCredentials: true,
	}))
	app.Use(helmet.New())
	app.Use(logger.New())

	MainRouter(app)
	if err = app.Listen(port); err != nil {
		log.Fatalf("Failed listening fiber app with port: %s", port)
	}
}
