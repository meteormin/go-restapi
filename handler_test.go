package go_restapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/app"
	"github.com/smyoo-pb/testclient/pkg/restapi"
	"io"
	"net/http/httptest"
	"testing"
)

var h restapi.Handler[TestEntity, *TestReq, *TestRes]

func init() {
	gofiber.New()
}

func TestNewHandler(t *testing.T) {
	h = restapi.NewHandler[TestEntity, *TestReq, *TestRes](
		&TestReq{},
		restapi.NewService[TestEntity, *TestReq, *TestRes](
			restapi.NewRepository[TestEntity](db, TestEntity{}),
			&TestRes{},
		),
	)
	if _, ok := h.(restapi.Handler[TestEntity, *TestReq, *TestRes]); ok {
		t.Log("ok!")
		return
	}

	t.Error("failed new handler")
}

func TestGenericHandler_Create(t *testing.T) {
	app.App().Fiber().Post("/tests", h.Create)
	testCases := makeFakeData(10)
	for _, testCase := range testCases {
		marshal, err := json.Marshal(&testCase)
		if err != nil {
			return
		}
		t.Logf("req body: %s", string(marshal))
		body := bytes.NewReader(marshal)
		req := httptest.NewRequest("POST", "/tests", body)
		req.Header.Set("Content-Type", "application/json")
		test, err := app.App().Test(req)
		if err != nil {
			t.Error(err)
		}

		all, err := io.ReadAll(test.Body)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(all))
	}
}

func TestGenericHandler_All(t *testing.T) {
	app.App().Fiber().Get("/tests", h.All)
	req := httptest.NewRequest("GET", "/tests", nil)
	test, err := app.App().Test(req)
	if err != nil {
		t.Error(err)
	}
	all, err := io.ReadAll(test.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(all))
}