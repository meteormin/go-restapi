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

type HandlerHook[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	Hook() HasHandlerEvent[Entity, Req, Res]
}

type ServiceHook[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	Hook() HasServiceEvent[Entity, Req, Res]
}

type HandlerEvents[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	ParseRequest(pr func(ctx *fiber.Ctx, dto Req) error)
	BeforeCallService(bs func(dto Req) error)
	AfterCallService(as func(dto Res) error)
}

type ServiceEvents[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	BeforeCallRepo(br func(dto Req, entity Entity) error)
	AfterCallRepo(ar func(dto Res, entity Entity) error)
}

type Features[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	parseRequest      *ParseRequest[Entity, Req]
	beforeCallService *BeforeCallService[Entity, Req, Res]
	afterCallService  *AfterCallService[Entity, Res]
	beforeCallRepo    *BeforeCallRepo[Entity, Req, Res]
	afterCallRepo     *AfterCallRepo[Entity, Res]
	methodEvent       MethodEvent
}

func (f *Features[Entity, Req, Res]) ParseRequest(pr func(ctx *fiber.Ctx, dto Req) error) {
	f.parseRequest = NewParseRequest[Entity, Req](pr)
}

func (f *Features[Entity, Req, Res]) BeforeCallService(bs func(dto Req) error) {
	f.beforeCallService = NewBeforeCallService[Entity, Req, Res](bs)
}

func (f *Features[Entity, Req, Res]) AfterCallService(as func(dto Res) error) {
	f.afterCallService = NewAfterCallService[Entity, Res](as)
}

func (f *Features[Entity, Req, Res]) BeforeCallRepo(br func(dto Req, entity Entity) error) {
	f.beforeCallRepo = NewBeforeCallRepo[Entity, Req, Res](br)
}

func (f *Features[Entity, Req, Res]) AfterCallRepo(ar func(dto Res, entity Entity) error) {
	f.afterCallRepo = NewAfterCallRepo[Entity, Res](ar)
}

type HasMethodEvent[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	methodEvent gollection.Collection[*Features[Entity, Req, Res]]
}

func (e *HasMethodEvent[Entity, Req, Res]) hasMethodEvent(ev MethodEvent) bool {
	return !e.methodEvent.Filter(func(v *Features[Entity, Req, Res], i int) bool {
		return ev == v.methodEvent
	}).IsEmpty()
}

func (e *HasMethodEvent[Entity, Req, Res]) getMethodEvent(ev MethodEvent) *Features[Entity, Req, Res] {
	if e.methodEvent == nil {
		e.methodEvent = gollection.NewCollection[*Features[Entity, Req, Res]](make([]*Features[Entity, Req, Res], 0))
	}
	first, _ := e.methodEvent.Filter(func(v *Features[Entity, Req, Res], i int) bool {
		return ev == v.methodEvent
	}).First()

	if first == nil {
		f := &Features[Entity, Req, Res]{}
		e.methodEvent.Add(f)
		return f
	}

	return *first
}

func (e *HasMethodEvent[Entity, Req, Res]) setMethodEvent(ev MethodEvent) {
	if e.methodEvent == nil {
		e.methodEvent = gollection.NewCollection[*Features[Entity, Req, Res]](make([]*Features[Entity, Req, Res], 0))
	}

	notExists := e.methodEvent.Filter(func(v *Features[Entity, Req, Res], i int) bool {
		return v.methodEvent == ev
	}).IsEmpty()

	if notExists {
		e.methodEvent.Add(&Features[Entity, Req, Res]{
			methodEvent: ev,
		})
	}
}

func (e *HasMethodEvent[Entity, Req, Res]) Create() *Features[Entity, Req, Res] {
	e.setMethodEvent(Create)
	return e.getMethodEvent(Create)
}

func (e *HasMethodEvent[Entity, Req, Res]) All() *Features[Entity, Req, Res] {
	e.setMethodEvent(All)
	return e.getMethodEvent(All)
}
func (e *HasMethodEvent[Entity, Req, Res]) Find() *Features[Entity, Req, Res] {
	e.setMethodEvent(Find)
	return e.getMethodEvent(Find)
}

func (e *HasMethodEvent[Entity, Req, Res]) Update() *Features[Entity, Req, Res] {
	e.setMethodEvent(Update)
	return e.getMethodEvent(Update)
}
func (e *HasMethodEvent[Entity, Req, Res]) Patch() *Features[Entity, Req, Res] {
	e.setMethodEvent(Patch)
	return e.getMethodEvent(Patch)
}
func (e *HasMethodEvent[Entity, Req, Res]) Delete() *Features[Entity, Req, Res] {
	e.setMethodEvent(Delete)
	return e.getMethodEvent(Delete)
}

type HasHandlerEvent[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	HasMethodEvent[Entity, Req, Res]
}

func (he *HasHandlerEvent[Entity, Req, Res]) Create() HandlerEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Create()
}
func (he *HasHandlerEvent[Entity, Req, Res]) Update() HandlerEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Update()
}
func (he *HasHandlerEvent[Entity, Req, Res]) Patch() HandlerEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Patch()
}
func (he *HasHandlerEvent[Entity, Req, Res]) Delete() HandlerEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Delete()
}
func (he *HasHandlerEvent[Entity, Req, Res]) Find() HandlerEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Find()
}
func (he *HasHandlerEvent[Entity, Req, Res]) All() HandlerEvents[Entity, Req, Res] {
	return he.HasMethodEvent.All()
}

type HasServiceEvent[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	HasMethodEvent[Entity, Req, Res]
}

func (he *HasServiceEvent[Entity, Req, Res]) Create() ServiceEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Create()
}
func (he *HasServiceEvent[Entity, Req, Res]) Update() ServiceEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Update()
}
func (he *HasServiceEvent[Entity, Req, Res]) Patch() ServiceEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Patch()
}
func (he *HasServiceEvent[Entity, Req, Res]) Delete() ServiceEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Delete()
}
func (he *HasServiceEvent[Entity, Req, Res]) Find() ServiceEvents[Entity, Req, Res] {
	return he.HasMethodEvent.Find()
}
func (he *HasServiceEvent[Entity, Req, Res]) All() ServiceEvents[Entity, Req, Res] {
	return he.HasMethodEvent.All()
}

type ParseRequest[Entity interface{}, Req RequestDTO[*Entity]] struct {
	event   Event
	handler func(ctx *fiber.Ctx, dto Req) error
}

func (pr *ParseRequest[Entity, Req]) Handler() func(ctx *fiber.Ctx, dto Req) error {
	return pr.handler
}

func NewParseRequest[Entity interface{}, Req RequestDTO[*Entity]](handler func(ctx *fiber.Ctx, dto Req) error) *ParseRequest[Entity, Req] {
	return &ParseRequest[Entity, Req]{
		event:   ParseRequestEvent,
		handler: handler,
	}
}

type BeforeCallService[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	event   Event
	handler func(dto Req) error
	HasMethodEvent[Entity, Req, Res]
}

func NewBeforeCallService[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](bs func(dto Req) error) *BeforeCallService[Entity, Req, Res] {
	return &BeforeCallService[Entity, Req, Res]{
		event:   BeforeCallServiceEvent,
		handler: bs,
	}
}
func (bs *BeforeCallService[Entity, Req, Res]) Handler() func(dto Req) error {
	return bs.handler
}

type AfterCallService[Entity interface{}, Res ResponseDTO[Entity]] struct {
	event   Event
	handler func(dto Res) error
}

func NewAfterCallService[Entity interface{}, Res ResponseDTO[Entity]](as func(dto Res) error) *AfterCallService[Entity, Res] {
	return &AfterCallService[Entity, Res]{
		event:   AfterCallServiceEvent,
		handler: as,
	}
}
func (as *AfterCallService[Entity, Res]) Handler() func(dto Res) error {
	return as.handler
}

type BeforeCallRepo[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	event   Event
	handler func(dto Req, entity Entity) error
}

func NewBeforeCallRepo[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](br func(dto Req, entity Entity) error) *BeforeCallRepo[Entity, Req, Res] {
	return &BeforeCallRepo[Entity, Req, Res]{
		event:   BeforeCallRepoEvent,
		handler: br,
	}
}
func (br *BeforeCallRepo[Entity, Req, Res]) Handler() func(dto Req, entity Entity) error {
	return br.handler
}

type AfterCallRepo[Entity interface{}, Res ResponseDTO[Entity]] struct {
	event   Event
	handler func(dto Res, entity Entity) error
}

func NewAfterCallRepo[Entity interface{}, Res ResponseDTO[Entity]](ar func(dto Res, entity Entity) error) *AfterCallRepo[Entity, Res] {
	return &AfterCallRepo[Entity, Res]{
		event:   AfterCallRepoEvent,
		handler: ar,
	}
}

func (ar *AfterCallRepo[Entity, Res]) Handler() func(dto Res, entity Entity) error {
	return ar.handler
}
