package main

import (
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/3-auth/handler"
	"github.com/Akihira77/gojobber/services/3-auth/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func NewHttpServer(db *gorm.DB, cld *util.Cloudinary, ccs *handler.GRPCClients) {
	port := os.Getenv("PORT")
	app := fiber.New(fiber.Config{
		BodyLimit:     5 * 1024 * 1024,
		CaseSensitive: true,
		StrictRouting: true,
		// Prefork:       true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("GATEWAY_URL"),
		AllowCredentials: true,
	}))
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(helmet.New())
	app.Use(logger.New())
	MainRouter(app, db, cld, ccs)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed listening fiber app with port: %s", port)
	}
}
