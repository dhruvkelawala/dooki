package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tokidooki/dooki/server/entities"
	er "github.com/tokidooki/dooki/server/errors"
)

func CreateMember(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "")
	if name == "" {
		return fiber.NewError(fiber.StatusNotAcceptable, er.NameIdUnspecified)
	}

	m, err := entities.GenerateMember(name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, er.ErrGenStruct)
	}
	return ctx.JSON(m)
}
