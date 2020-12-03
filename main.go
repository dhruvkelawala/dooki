package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Static("/", "./client/build/")

	app.Static("/static/", "./client/build/static/")

	app.Get("/api/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
