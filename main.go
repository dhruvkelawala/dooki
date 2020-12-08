package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)


func getPort() string {
	port, ok := os.LookupEnv("PORT")
	if ok {
		return fmt.Sprintf(":%s", port)
	}

	return ":3000"
}

func main() {
	router:= fiber.New()

	router.Use(cors.New())

	//Static
	router.Static("/", "./client/build/")
	router.Static("/static/", "./client/build/static/")

	router.Get("/api/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(router.Listen(getPort()))
}
