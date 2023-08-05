package restapi_test

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/gofiber/app"
	"io"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	route := restapi.New[TestEntity, *TestReq, *TestRes](
		db,
		TestEntity{},
		&TestReq{},
		&TestRes{},
	)

	app.App().Route("/api", func(router app.Router, app app.Application) {
		router.Route("/test", route)
	})
	app.App().Bootstrap()
}

func TestRoute(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test", nil)
	test, err := app.App().Fiber().Test(req)
	if err != nil {
		t.Error(err)
	}

	all, err := io.ReadAll(test.Body)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(all))
}
