package restapi

import "github.com/gofiber/fiber/v2"

type Event string

type Features struct {
	ParseRequest      *ParseRequest
	BeforeCallService *BeforeCallService
	AfterCallService  *AfterCallService
	BeforeCallRepo    *BeforeCallRepo
	AfterCallRepo     *AfterCallRepo
}

const (
	ParseRequestEvent      Event = "parseRequest"
	BeforeCallServiceEvent Event = "beforeCallService"
	AfterCallServiceEvent  Event = "afterCallService"
	BeforeCallRepoEvent    Event = "beforeCallRepo"
	AfterCallRepoEvent     Event = "afterCallRepo"
)

type HandlerEvents interface {
	ParseRequest(pr ParseRequestHandler)
	BeforeCallService(bs BeforeCallServiceHandler)
	AfterCallService(as AfterCallServiceHandler)
}

type ServiceEvents interface {
	BeforeCallRepo(br BeforeCallRepoHandler)
	AfterCallRepo(ar AfterCallRepoHandler)
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

type BeforeCallServiceHandler = func(dto interface{}) error
type BeforeCallService struct {
	event   Event
	handler BeforeCallServiceHandler
}

func NewBeforeCallService(bs BeforeCallServiceHandler) *BeforeCallService {
	return &BeforeCallService{
		event:   BeforeCallServiceEvent,
		handler: bs,
	}
}
func (bs *BeforeCallService) Handler() BeforeCallServiceHandler {
	return bs.handler
}

type AfterCallServiceHandler = func(dto interface{}) error
type AfterCallService struct {
	event   Event
	handler AfterCallServiceHandler
}

func NewAfterCallService(as AfterCallServiceHandler) *AfterCallService {
	return &AfterCallService{
		event:   AfterCallServiceEvent,
		handler: as,
	}
}
func (as *AfterCallService) Handler() AfterCallServiceHandler {
	return as.handler
}

type BeforeCallRepoHandler = func(dto interface{}, entity interface{}) error
type BeforeCallRepo struct {
	event   Event
	handler BeforeCallRepoHandler
}

func NewBeforeCallRepo(br BeforeCallRepoHandler) *BeforeCallRepo {
	return &BeforeCallRepo{
		event:   BeforeCallRepoEvent,
		handler: br,
	}
}
func (br *BeforeCallRepo) Handler() BeforeCallRepoHandler {
	return br.handler
}

type AfterCallRepoHandler = func(dto interface{}, entity interface{}) error
type AfterCallRepo struct {
	event   Event
	handler AfterCallRepoHandler
}

func NewAfterCallRepo(ar AfterCallRepoHandler) *AfterCallRepo {
	return &AfterCallRepo{
		event:   AfterCallRepoEvent,
		handler: ar,
	}
}
func (ar *AfterCallRepo) Handler() AfterCallRepoHandler {
	return ar.handler
}
