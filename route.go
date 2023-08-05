package restapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
)

func Route[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](handler Handler[Entity, Req, Res]) app.SubRouter {
	return func(router fiber.Router) {
		router.Post("/", handler.Create)
		router.Get("/", handler.All)
		router.Get("/:id", handler.Find)
		router.Put("/:id", handler.Update)
		router.Patch("/:id", handler.Patch)
		router.Delete("/:id", handler.Delete)
	}
}
