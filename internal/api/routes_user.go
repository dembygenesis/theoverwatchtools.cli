package api

import "github.com/gofiber/fiber/v2"

func (a *Api) DummyUser(ctx *fiber.Ctx) error {
	return ctx.SendString("Dummy user routes")
}
