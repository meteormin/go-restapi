package go_restapi_test

import (
	"github.com/smyoo-pb/testclient/pkg/restapi"
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
