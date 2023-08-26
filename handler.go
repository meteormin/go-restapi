package restapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/utils"
	"log"
	"strconv"
)

type Handler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	All(ctx *fiber.Ctx) error
	Find(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Patch(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	GetService() Service[Entity, Req, Res]
	HandlerHook[Entity, Req, Res]
}

type GenericHandler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	req     Req
	service Service[Entity, Req, Res]
	events  *HasHandlerEvent[Entity, Req, Res]
}

func NewHandler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](req Req, service Service[Entity, Req, Res]) Handler[Entity, Req, Res] {
	return &GenericHandler[Entity, Req, Res]{
		req:     req,
		service: service,
		events: &HasHandlerEvent[Entity, Req, Res]{
			&HasMethodEvent[Entity, Req, Res]{methodEvent: nil},
		},
	}
}

func (g *GenericHandler[Entity, Req, Res]) All(ctx *fiber.Ctx) error {
	if g.features(All).parseRequest != nil {
		err := g.features(All).parseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	filter, err := NewFilter[Entity](ctx, g.req)

	if g.features(All).beforeCallService != nil {
		err = g.features(All).beforeCallService.handler(g.req)
		if err != nil {
			return err
		}
	}

	all, err := g.service.All(filter)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(all)
}

func (g *GenericHandler[Entity, Req, Res]) Find(ctx *fiber.Ctx) error {
	if g.features(Find).parseRequest != nil {
		err := g.features(Find).parseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	params := ctx.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	find, err := g.service.Find(uint(pk))
	if err != nil {
		return err
	}

	log.Print(g.features(Find))
	if g.features(Find).afterCallService != nil {
		err = g.features(Find).afterCallService.handler(find)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(find)
}

func (g *GenericHandler[Entity, Req, Res]) Create(ctx *fiber.Ctx) error {
	if g.features(Create).parseRequest != nil {
		err := g.features(Create).parseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	req := g.req
	errRes := utils.HandleValidate(ctx, req)
	if errRes != nil {
		return errRes.Response()
	}

	if g.features(Create).beforeCallService != nil {
		err := g.features(Create).beforeCallService.handler(req)
		if err != nil {
			return err
		}
	}

	create, err := g.service.Create(req)
	if err != nil {
		return err
	}

	if g.features(Create).afterCallService != nil {
		err = g.features(Create).afterCallService.handler(create)
		if err != nil {
			return err
		}
	}
	return ctx.Status(fiber.StatusCreated).JSON(create)
}

func (g *GenericHandler[Entity, Req, Res]) Update(ctx *fiber.Ctx) error {
	if g.features(Update).parseRequest != nil {
		err := g.features(Update).parseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	params := ctx.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	req := g.req
	errRes := utils.HandleValidate(ctx, req)
	if errRes != nil {
		return errRes.Response()
	}

	if g.features(Update).beforeCallService != nil {
		err = g.features(Update).beforeCallService.handler(req)
		if err != nil {
			return err
		}
	}

	update, err := g.service.Update(uint(pk), req)
	if err != nil {
		return err
	}

	if g.features(Update).afterCallService != nil {
		err = g.features(Update).afterCallService.handler(update)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(update)
}

func (g *GenericHandler[Entity, Req, Res]) Patch(ctx *fiber.Ctx) error {
	if g.features(Patch).parseRequest != nil {
		err := g.features(Patch).parseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	params := ctx.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	req := g.req
	errRes := utils.HandleValidate(ctx, req)
	if errRes != nil {
		return errRes.Response()
	}

	if g.features(Patch).beforeCallService != nil {
		err = g.features(Patch).beforeCallService.handler(req)
		if err != nil {
			return err
		}
	}

	update, err := g.service.Patch(uint(pk), req)
	if err != nil {
		return err
	}

	if g.features(Patch).afterCallService != nil {
		err = g.features(Patch).afterCallService.handler(update)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(update)
}

func (g *GenericHandler[Entity, Req, Res]) Delete(ctx *fiber.Ctx) error {
	if g.features(Delete).parseRequest != nil {
		err := g.features(Delete).parseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	params := ctx.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	find, err := g.service.Delete(uint(pk))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(map[string]interface{}{
		"result": find,
	})
}

func (g *GenericHandler[Entity, Req, Res]) GetService() Service[Entity, Req, Res] {
	return g.service
}

func (g *GenericHandler[Entity, Req, Res]) Hook() *HasHandlerEvent[Entity, Req, Res] {
	return g.events
}

func (g *GenericHandler[Entity, Req, Res]) features(event MethodEvent) *Features[Entity, Req, Res] {
	return g.events.getMethodEvent(event)
}
