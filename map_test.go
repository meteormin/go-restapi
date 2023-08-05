package go_restapi_test

import (
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/structs"
	"github.com/smyoo-pb/testclient/pkg/restapi"
	"reflect"
	"testing"
)

func init() {
	log.GetLogger().Debug("test")
}

type DestStruct struct {
	Id       *uint
	Name     string
	IsBool   bool
	PtrFloat float64
}

type SrcStruct struct {
	Id       uint
	Name     string
	IsBool   bool
	PtrFloat *float64
}

func TestMap(t *testing.T) {
	dest := &DestStruct{}

	ptrFloat := 1.0
	src := &SrcStruct{
		Id:       1,
		Name:     "src2",
		IsBool:   true,
		PtrFloat: &ptrFloat,
	}

	err := restapi.Map(dest, src)
	if err != nil {
		t.Error(err)
	}
	t.Logf("dest > %v", dest)
	t.Logf("%d", *dest.Id)
}

func TestStructs(t *testing.T) {
	var id uint = 1
	dest := DestStruct{Id: &id, Name: "name"}
	destStruct := structs.New(&dest)
	t.Logf("%v", destStruct.Fields()[0].Kind())

	v := destStruct.Fields()[0].Value()
	ty := reflect.TypeOf(v)
	t.Logf("%v", ty.Elem().Kind())
	if ty.Elem().Kind() == reflect.Uint {
		i, _ := v.(*uint)
		t.Logf("%v", *i)
	}
}
