package restapi_test

import (
	"github.com/go-faker/faker/v4"
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/gofiber/database"
	"gorm.io/gorm"
)

type TestEntity struct {
	gorm.Model
	Name              string `gorm:"column:name;unique"`
	TestRelationModel TestRelationModel
}

type TestRelationModel struct {
	gorm.Model
	TestEntityId uint `gorm:"column:test_entity_id"`
	Seq          int
}

type TestRelationReq struct {
	TestEntityId uint `faker:"-" json:"test_entity_id"`
	Seq          int  `json:"seq" validate:"required"`
}

type TestReq struct {
	restapi.RequestDTO[TestEntity] `faker:"-" json:"-"`
	Name                           string          `faker:"username" json:"name" validate:"required"`
	TestRelation                   TestRelationReq `json:"test_relation" validate:"required"`
}

func (tr *TestReq) ToEntity(ent *TestEntity) error {
	ent.Name = tr.Name
	ent.TestRelationModel = TestRelationModel{
		Seq: tr.TestRelation.Seq,
	}
	return nil
}

type TestRelationRes struct {
	Id           uint `json:"id"`
	TestEntityId uint `json:"test_entity_id"`
	Seq          int  `json:"seq"`
}

type TestRes struct {
	restapi.ResponseDTO[TestEntity] `json:"-"`
	Id                              uint            `json:"id"`
	Name                            string          `json:"name"`
	TestRelationRes                 TestRelationRes `json:"test_relation_res"`
}

func (tr *TestRes) FromEntity(ent TestEntity) error {
	tr.Id = ent.ID
	tr.Name = ent.Name
	tr.TestRelationRes = TestRelationRes{
		Id:           ent.TestRelationModel.ID,
		TestEntityId: ent.TestRelationModel.TestEntityId,
		Seq:          ent.TestRelationModel.Seq,
	}

	return nil
}

var db *gorm.DB

func init() {
	db = database.New(database.Config{
		Name:   "sqlite",
		Driver: "sqlite",
		Dbname: "test",
	})
	err := db.AutoMigrate(TestEntity{}, TestRelationModel{})
	if err != nil {
		panic(err)
	}

	db.Debug().Exec("DELETE FROM `test_entities`;")
	db.Debug().Exec("UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'test_entities';")
	db.Debug().Exec("DELETE FROM `test_relation_models`;")
	db.Debug().Exec("UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'test_relation_models';")
}

func makeFakeData(cnt int) []TestReq {
	testData := make([]TestReq, 0)
	for i := 0; i < cnt; i++ {
		data := TestReq{}

		data.TestRelation = TestRelationReq{
			Seq: i,
		}
		err := faker.FakeData(&data)
		if err != nil {
			panic(err)
		}
		testData = append(testData, data)
	}

	return testData
}
