package main

import (
	"fmt"
	"log"

	"jotham/database"
	"jotham/handler"
	"jotham/helper"
	"jotham/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if err := helper.MongoConnect(); err != nil {
		log.Fatal(err)
	}
	if err := database.PGConnect(); err != nil {
		fmt.Println(err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    52428800, //50mb
		ServerHeader: "Fiber",
		Prefork:      true,
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
	root := utils.GetROOTDomain()
	fmt.Printf("Root Domain set is %v", root)
	log.Fatal(app.Listen(":3000"))
}
