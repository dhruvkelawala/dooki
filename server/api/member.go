package api

import (
	"github.com/daemon1024/dokidoki/server/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateMember(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "")
	if name == "" {
		return fiber.NewError(fiber.StatusNotAcceptable, "name not specified")
	}

	id, _ := uuid.NewRandom()
	m := entities.Member{
		ID: id.String(),
		Name: name,
	}

	return ctx.JSON(m)
}
