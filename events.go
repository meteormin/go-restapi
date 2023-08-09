package restapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gollection"
)

type MethodEvent string

const (
	Common MethodEvent = "ALL"
	Create MethodEvent = "C"
	All    MethodEvent = "RA"
	Find   MethodEvent = "R"
	Update MethodEvent = "UA"
	Patch  MethodEvent = "U"
	Delete MethodEvent = "D"
)

type Event string

const (
	ParseRequestEvent      Event = "parseRequest"
	BeforeCallServiceEvent Event = "beforeCallService"
	AfterCallServiceEvent  Event = "afterCallService"
	BeforeCallRepoEvent    Event = "beforeCallRepo"
	AfterCallRepoEvent     Event = "afterCallRepo"
)

type Features struct {
	ParseRequest      *ParseRequest
	BeforeCallService *BeforeCallService
	AfterCallService  *AfterCallService
	BeforeCallRepo    *BeforeCallRepo
	AfterCallRepo     *AfterCallRepo
	methodEvent       MethodEvent
}

type HasMethodEvent struct {
	methodEvent gollection.Collection[Features]
}

func (e *HasMethodEvent) hasMethodEvent(ev MethodEvent) bool {
	return !e.methodEvent.Filter(func(v Features, i int) bool {
		return ev == v.methodEvent
	}).IsEmpty()
}

func (e *HasMethodEvent) getMethodEvent(ev MethodEvent) (*Features, error) {
	return e.methodEvent.Filter(func(v Features, i int) bool {
		return ev == v.methodEvent
	}).First()
}

func (e *HasMethodEvent) setMethodEvent(ev MethodEvent) {
	if e.methodEvent == nil {
		e.methodEvent = gollection.NewCollection[Features](make([]Features, 0))
	}

	notExists := e.methodEvent.Filter(func(v Features, i int) bool {
		return v.methodEvent == ev
	}).IsEmpty()

	if notExists {
		e.methodEvent.Add(Features{
			methodEvent: ev,
		})
	}
}

func (e *HasMethodEvent) Create() (*Features, error) {
	e.setMethodEvent(Create)
	return e.getMethodEvent(Create)
}

func (e *HasMethodEvent) All() (*Features, error) {
	e.setMethodEvent(All)
	return e.getMethodEvent(All)
}
func (e *HasMethodEvent) Find() (*Features, error) {
	e.setMethodEvent(Find)
	return e.getMethodEvent(Find)
}

func (e *HasMethodEvent) Update() (*Features, error) {
	e.setMethodEvent(Update)
	return e.getMethodEvent(Update)
}
func (e *HasMethodEvent) Patch() (*Features, error) {
	e.setMethodEvent(Patch)
	return e.getMethodEvent(Patch)
}
func (e *HasMethodEvent) Delete() (*Features, error) {
	e.setMethodEvent(Delete)
	return e.getMethodEvent(Delete)
}

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
	HasMethodEvent
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
