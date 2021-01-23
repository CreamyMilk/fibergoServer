package handler

import "github.com/gofiber/fiber/v2"

func NofileHandler(c *fiber.Ctx) error {
	c.SendFile("./public/404.html")
	return c.Status(404).JSON(&fiber.Map{
		"success": false,
		"error":   "There are no posts!",
	})
}
