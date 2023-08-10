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

type HandlerHook interface {
	Hook() HasHandlerEvent
}

type ServiceHook interface {
	Hook() HasServiceEvent
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

type Features struct {
	parseRequest      *ParseRequest
	beforeCallService *BeforeCallService
	afterCallService  *AfterCallService
	beforeCallRepo    *BeforeCallRepo
	afterCallRepo     *AfterCallRepo
	methodEvent       MethodEvent
}

func (f *Features) ParseRequest(pr ParseRequestHandler) {
	f.parseRequest = NewParseRequest(pr)
}

func (f *Features) BeforeCallService(bs BeforeCallServiceHandler) {
	f.beforeCallService = NewBeforeCallService(bs)
}

func (f *Features) AfterCallService(as AfterCallServiceHandler) {
	f.afterCallService = NewAfterCallService(as)
}

func (f *Features) BeforeCallRepo(br BeforeCallRepoHandler) {
	f.beforeCallRepo = NewBeforeCallRepo(br)
}

func (f *Features) AfterCallRepo(ar AfterCallRepoHandler) {
	f.afterCallRepo = NewAfterCallRepo(ar)
}

type HasMethodEvent struct {
	methodEvent gollection.Collection[Features]
}

func (e *HasMethodEvent) hasMethodEvent(ev MethodEvent) bool {
	return !e.methodEvent.Filter(func(v Features, i int) bool {
		return ev == v.methodEvent
	}).IsEmpty()
}

func (e *HasMethodEvent) getMethodEvent(ev MethodEvent) *Features {
	if e.methodEvent == nil {
		e.methodEvent = gollection.NewCollection[Features](make([]Features, 0))
	}
	first, _ := e.methodEvent.Filter(func(v Features, i int) bool {
		return ev == v.methodEvent
	}).First()

	if first == nil {
		return &Features{}
	}

	return first
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

func (e *HasMethodEvent) Create() *Features {
	e.setMethodEvent(Create)
	return e.getMethodEvent(Create)
}

func (e *HasMethodEvent) All() *Features {
	e.setMethodEvent(All)
	return e.getMethodEvent(All)
}
func (e *HasMethodEvent) Find() *Features {
	e.setMethodEvent(Find)
	return e.getMethodEvent(Find)
}

func (e *HasMethodEvent) Update() *Features {
	e.setMethodEvent(Update)
	return e.getMethodEvent(Update)
}
func (e *HasMethodEvent) Patch() *Features {
	e.setMethodEvent(Patch)
	return e.getMethodEvent(Patch)
}
func (e *HasMethodEvent) Delete() *Features {
	e.setMethodEvent(Delete)
	return e.getMethodEvent(Delete)
}

type HasHandlerEvent struct {
	HasMethodEvent
}

func (he *HasHandlerEvent) Create() HandlerEvents {
	return he.HasMethodEvent.Create()
}
func (he *HasHandlerEvent) Update() HandlerEvents {
	return he.HasMethodEvent.Update()
}
func (he *HasHandlerEvent) Patch() HandlerEvents {
	return he.HasMethodEvent.Patch()
}
func (he *HasHandlerEvent) Delete() HandlerEvents {
	return he.HasMethodEvent.Delete()
}
func (he *HasHandlerEvent) Find() HandlerEvents {
	return he.HasMethodEvent.Find()
}
func (he *HasHandlerEvent) All() HandlerEvents {
	return he.HasMethodEvent.All()
}

type HasServiceEvent struct {
	HasMethodEvent
}

func (he *HasServiceEvent) Create() ServiceEvents {
	return he.HasMethodEvent.Create()
}
func (he *HasServiceEvent) Update() ServiceEvents {
	return he.HasMethodEvent.Update()
}
func (he *HasServiceEvent) Patch() ServiceEvents {
	return he.HasMethodEvent.Patch()
}
func (he *HasServiceEvent) Delete() ServiceEvents {
	return he.HasMethodEvent.Delete()
}
func (he *HasServiceEvent) Find() ServiceEvents {
	return he.HasMethodEvent.Find()
}
func (he *HasServiceEvent) All() ServiceEvents {
	return he.HasMethodEvent.All()
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
