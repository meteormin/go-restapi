package restapi

import (
	"github.com/gofiber/fiber/v2"
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
	HandlerEvents
}

type GenericHandler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	req      Req
	service  Service[Entity, Req, Res]
	features Features
}

func NewHandler[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](req Req, service Service[Entity, Req, Res]) Handler[Entity, Req, Res] {
	return &GenericHandler[Entity, Req, Res]{
		req:     req,
		service: service,
		features: Features{
			ParseRequest:     nil,
			BeforeValidation: nil,
			AfterValidation:  nil,
		},
	}
}

func (g *GenericHandler[Entity, Req, Res]) All(ctx *fiber.Ctx) error {
	if g.features.ParseRequest != nil {
		err := g.features.ParseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	filter, err := NewFilter[Entity](ctx, g.req)
	all, err := g.service.All(filter)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(all)
}

func (g *GenericHandler[Entity, Req, Res]) Find(ctx *fiber.Ctx) error {
	if g.features.ParseRequest != nil {
		err := g.features.ParseRequest.handler(ctx, g.req)
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

	return ctx.Status(fiber.StatusOK).JSON(find)
}

func (g *GenericHandler[Entity, Req, Res]) Create(ctx *fiber.Ctx) error {
	if g.features.ParseRequest != nil {
		err := g.features.ParseRequest.handler(ctx, g.req)
		if err != nil {
			return err
		}
	}

	req := g.req
	errRes := utils.HandleValidate(ctx, req)
	if errRes != nil {
		return errRes.Response()
	}

	create, err := g.service.Create(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(create)
}

func (g *GenericHandler[Entity, Req, Res]) Update(ctx *fiber.Ctx) error {
	if g.features.ParseRequest != nil {
		err := g.features.ParseRequest.handler(ctx, g.req)
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

	update, err := g.service.Update(uint(pk), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(update)
}

func (g *GenericHandler[Entity, Req, Res]) Patch(ctx *fiber.Ctx) error {
	if g.features.ParseRequest != nil {
		err := g.features.ParseRequest.handler(ctx, g.req)
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

	update, err := g.service.Patch(uint(pk), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(update)
}

func (g *GenericHandler[Entity, Req, Res]) Delete(ctx *fiber.Ctx) error {
	if g.features.ParseRequest != nil {
		err := g.features.ParseRequest.handler(ctx, g.req)
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

func (g *GenericHandler[Entity, Req, Res]) ParseRequest(pr ParseRequestHandler) {
	g.features.ParseRequest = NewParseRequest(pr)
}

func (g *GenericHandler[Entity, Req, Res]) BeforeValidation(bv BeforeValidationHandler) {
	g.features.BeforeValidation = NewBeforeValidation(bv)
}

func (g *GenericHandler[Entity, Req, Res]) AfterValidation(av AfterValidationHandler) {
	g.features.AfterValidation = NewAfterValidation(av)
}
