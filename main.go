package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"jotham/database"
	"jotham/handler"
	"jotham/helper"
)

func main() {
	if err := helper.MongoConnect(); err != nil {
		log.Fatal(err)
	}
	// Connect with database
	if err := database.PGConnect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    52428800, //50mb
		ServerHeader: "Fiber",
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/savepdf", handler.SaveHandler)
	app.Static("/", "./public")
	app.Static("/pages", "./uploads/testfolder", fiber.Static{
		Compress:  false,
		ByteRange: false,
		Browse:    true,
	})

	app.Use(handler.NofileHandler)
	log.Fatal(app.Listen(":3000"))
}
