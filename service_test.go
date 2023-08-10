package restapi_test

import (
	"github.com/miniyus/go-restapi"
	"testing"
)

var serviceTestData []TestReq

func init() {
	serviceTestData = makeFakeData(10)
}

var service restapi.Service[TestEntity, *TestReq, *TestRes]

func TestNewService(t *testing.T) {
	repo := restapi.NewRepository[TestEntity](db, TestEntity{})
	service = restapi.NewService[TestEntity, *TestReq, *TestRes](
		repo,
		&TestRes{},
	)
}

func TestGenericService_Create(t *testing.T) {
	entities := make([]TestRes, 0)
	for _, data := range serviceTestData {
		create, err := service.Create(&data)
		if err != nil {
			t.Error(err)
		}
		entities = append(entities, *create)
	}

	t.Log(entities)
}

func TestGenericService_All(t *testing.T) {
	filter := restapi.Filter[TestEntity]{}
	all, err := service.All(&filter)
	if err != nil {
		t.Error(err)
	}

	for _, res := range all {
		t.Logf("%v", res)
	}
}

func TestGenericService_Hook(t *testing.T) {
	hook := service.Hook()
	find := hook.Create()
	find.AfterCallRepo(func(dto interface{}, entity interface{}) error {
		t.Log(dto)
		t.Log(entity)
		return nil
	})

	res, err := service.Find(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
