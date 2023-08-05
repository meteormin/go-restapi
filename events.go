package restapi

import "github.com/gofiber/fiber/v2"

type Event string

type Features struct {
	ParseRequest     *ParseRequest
	BeforeValidation *BeforeValidation
	AfterValidation  *AfterValidation
}

const (
	ParseRequestEvent      Event = "parseRequest"
	BeforeValidationEvent  Event = "beforeValidation"
	AfterValidationEvent   Event = "afterValidation"
	BeforeCallServiceEvent Event = "beforeCallService"
	AfterCallServiceEvent  Event = "afterCallService"
	BeforeCallRepoEvent    Event = "beforeCallRepo"
	AfterCallRepoEvent     Event = "afterCallRepo"
	BeforeResponse         Event = "beforeResponse"
	CreatedResponse        Event = "createdResponse"
)

type HandlerEvents interface {
	ParseRequest(pr ParseRequestHandler)
	BeforeValidation(bv BeforeValidationHandler)
	AfterValidation(av AfterValidationHandler)
}

type ParseRequestHandler = func(ctx *fiber.Ctx, dto interface{}) error

type ParseRequest struct {
	event   Event
	handler ParseRequestHandler
}

func (pr *ParseRequest) Handler() ParseRequestHandler {
	return pr.handler
}

func NewParseRequest(handler ParseRequestHandler) *ParseRequest {
	return &ParseRequest{
		event:   ParseRequestEvent,
		handler: handler,
	}
}

type BeforeValidationHandler = func(dto interface{}) error
type BeforeValidation struct {
	event   Event
	handler BeforeValidationHandler
}

func NewBeforeValidation(handler BeforeValidationHandler) *BeforeValidation {
	return &BeforeValidation{
		event:   BeforeValidationEvent,
		handler: handler,
	}
}
func (bv *BeforeValidation) Handler() func(dto interface{}) error {
	return bv.handler
}

type AfterValidationHandler = func(dto interface{}) error
type AfterValidation struct {
	event   Event
	handler AfterValidationHandler
}

func NewAfterValidation(handler AfterValidationHandler) *AfterValidation {
	return &AfterValidation{
		event:   AfterValidationEvent,
		handler: handler,
	}
}
func (av *AfterValidation) Handler() AfterValidationHandler {
	return av.handler
}
