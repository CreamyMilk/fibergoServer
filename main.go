package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/crypto/acme/autocert"

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
		fmt.Println(err)
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

	// Let’s Encrypt has rate limits: https://letsencrypt.org/docs/rate-limits/
	// It's recommended to use it's staging environment to test the code:
	// https://letsencrypt.org/docs/staging-environment/

	// Certificate manager
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Replace with your domain
		HostPolicy: autocert.HostWhitelist("readzy.africa"),
		// Folder to store the certificates
		Cache: autocert.DirCache("./certs"),
	}

	// TLS Config
	cfg := &tls.Config{
		// Get Certificate from Let's Encrypt
		GetCertificate: m.GetCertificate,
		// By default NextProtos contains the "h2"
		// This has to be removed since Fasthttp does not support HTTP/2
		// Or it will cause a flood of PRI method logs
		// http://webconcepts.info/concepts/http-method/PRI
		NextProtos: []string{
			"http/1.1", "acme-tls/1",
		},
	}
	ln, err := tls.Listen("tcp", ":3000", cfg)
	if err != nil {
		panic(err)
	}

	// Start server
	log.Fatal(app.Listener(ln))
	log.Fatal(app.Listen(":3001"))
}
