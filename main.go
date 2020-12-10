package main

import (
	"fmt"
	"github.com/daemon1024/dokidoki/server/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
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

	//Middleware
	router.Use(cors.New())
	router.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			return ctx.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	//Static
	router.Static("/", "./client/build/")
	router.Static("/static/", "./client/build/static/")

	//REST
	router.Post("/room/:id/modify", api.ModifyRoomStatus)
	router.Post("/room/create", api.CreateRoom)
	router.Post("/member/create", api.CreateMember)

	//Websocket
	router.Get("/ws/room/:id", websocket.New(api.MyRoom))

	log.Fatal(router.Listen(getPort()))
}
