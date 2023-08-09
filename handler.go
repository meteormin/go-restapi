package restapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/utils"
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
}

type GenericHandler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	req     Req
	service Service[Entity, Req, Res]
	events  HasMethodEvent
}

func NewHandler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](req Req, service Service[Entity, Req, Res]) Handler[Entity, Req, Res] {
	return &GenericHandler[Entity, Req, Res]{
		req:     req,
		service: service,
		events: HasMethodEvent{
			methodEvent: nil,
		},
	}
}

func (g *GenericHandler[Entity, Req, Res]) All(ctx *fiber.Ctx) error {
	if g.features(All).ParseRequest != nil {
		err := g.features(All).ParseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	filter, err := NewFilter[Entity](ctx, g.req)

	if g.features(All).BeforeCallService != nil {
		err = g.features(All).BeforeCallService.handler(g.req)
		if err != nil {
			return err
		}
	}

	all, err := g.service.All(filter)

	if err != nil {
		return err
	}

	if g.features(All).AfterCallService != nil {
		err = g.features(All).AfterCallService.handler(all)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(all)
}

func (g *GenericHandler[Entity, Req, Res]) Find(ctx *fiber.Ctx) error {
	if g.features(Find).ParseRequest != nil {
		err := g.features(Find).ParseRequest.handler(ctx, g.req)
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

	if g.features(Find).AfterCallService != nil {
		err = g.features(Find).AfterCallService.handler(find)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(find)
}

func (g *GenericHandler[Entity, Req, Res]) Create(ctx *fiber.Ctx) error {
	if g.features(Create).ParseRequest != nil {
		err := g.features(Create).ParseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	req := g.req
	errRes := utils.HandleValidate(ctx, req)
	if errRes != nil {
		return errRes.Response()
	}

	if g.features(Create).BeforeCallService != nil {
		err := g.features(Create).BeforeCallService.handler(req)
		if err != nil {
			return err
		}
	}

	create, err := g.service.Create(req)
	if err != nil {
		return err
	}

	if g.features(Create).AfterCallService != nil {
		err = g.features(Create).AfterCallService.handler(req)
		if err != nil {
			return err
		}
	}
	return ctx.Status(fiber.StatusCreated).JSON(create)
}

func (g *GenericHandler[Entity, Req, Res]) Update(ctx *fiber.Ctx) error {
	if g.features(Update).ParseRequest != nil {
		err := g.features(Update).ParseRequest.handler(ctx, g.req)
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

	if g.features(Update).BeforeCallService != nil {
		err = g.features(Update).BeforeCallService.handler(req)
		if err != nil {
			return err
		}
	}

	update, err := g.service.Update(uint(pk), req)
	if err != nil {
		return err
	}

	if g.features(Update).AfterCallService != nil {
		err = g.features(Update).AfterCallService.handler(update)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(update)
}

func (g *GenericHandler[Entity, Req, Res]) Patch(ctx *fiber.Ctx) error {
	if g.features(Patch).ParseRequest != nil {
		err := g.features(Patch).ParseRequest.handler(ctx, g.req)
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

	if g.features(Patch).BeforeCallService != nil {
		err = g.features(Patch).BeforeCallService.handler(req)
		if err != nil {
			return err
		}
	}

	update, err := g.service.Patch(uint(pk), req)
	if err != nil {
		return err
	}

	if g.features(Patch).AfterCallService != nil {
		err = g.features(Patch).AfterCallService.handler(update)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(update)
}

func (g *GenericHandler[Entity, Req, Res]) Delete(ctx *fiber.Ctx) error {
	if g.features(Delete).ParseRequest != nil {
		err := g.features(Delete).ParseRequest.handler(ctx, g.req)
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

	if g.features(Delete).AfterCallService != nil {
		err = g.features(Delete).AfterCallService.handler(find)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusNoContent).JSON(map[string]interface{}{
		"result": find,
	})
}

func (g *GenericHandler[Entity, Req, Res]) GetService() Service[Entity, Req, Res] {
	return g.service
}

func (g *GenericHandler[Entity, Req, Res]) Hook() HasMethodEvent {
	return g.events
}

func (g *GenericHandler[Entity, Req, Res]) features(event MethodEvent) *Features {
	feat, err := g.events.getMethodEvent(event)
	gofiber.Log().Error(err)

	return feat
}
