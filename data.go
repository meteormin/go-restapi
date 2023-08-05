package restapi

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/structs"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type RequestDTO[Entity interface{}] interface {
	ToEntity(ent Entity) error
}

type ResponseDTO[Entity interface{}] interface {
	FromEntity(ent Entity) error
}

type Sorter struct {
	columns []string
}

func (s *Sorter) Sort(db *gorm.DB) *gorm.DB {
	for _, col := range s.columns {
		order := strings.Split(col, " ")
		if len(order) != 0 {
			db = db.Order(fmt.Sprintf("%s %s", order[0], order[1]))
		}
	}
	return db
}

type Searcher[Entity interface{}] struct {
	req    RequestDTO[*Entity]
	entity Entity
}

func (s *Searcher[Entity]) SetEntity(ent Entity) {
	s.entity = ent
}

func (s *Searcher[Entity]) Search(db *gorm.DB) (*gorm.DB, error) {
	if s.req == nil {
		return db, nil
	}

	err := s.req.ToEntity(&s.entity)
	if err != nil {
		return db, err
	}

	db.Where(s.entity)

	return db, nil
}

type Filter[Entity interface{}] struct {
	pagination.Page
	Sorter
	Searcher[Entity]
}

func NewFilter[Entity interface{}](ctx *fiber.Ctx, req RequestDTO[*Entity]) (*Filter[Entity], error) {
	param := ctx.Query("sort")
	sort := strings.Split(param, ",")
	fromCtx, _ := pagination.GetPageFromCtx(ctx)
	err := ctx.QueryParser(&req)
	if err != nil {
		return nil, err
	}

	return &Filter[Entity]{
		Page:   fromCtx,
		Sorter: Sorter{columns: sort},
		Searcher: Searcher[Entity]{
			req: req,
		},
	}, nil
}

func Map(dest, src interface{}) error {
	isStructDes := structs.IsStruct(dest)
	isStructSrc := structs.IsStruct(src)
	if !isStructDes || !isStructSrc {
		return errors.New("can't Map")
	}

	destFs := structs.Fields(dest)
	srcFs := structs.Fields(src)

	for _, srcF := range srcFs {
		for _, destF := range destFs {
			if srcF.Name() == destF.Name() {
				if srcF.Kind() != destF.Kind() {
					if srcF.Kind() == reflect.Ptr {
						val := transFromPtr(srcF.Value())
						err := destF.Set(val)
						if err != nil {
							log.GetLogger().Error(err)
							return fmt.Errorf("src ptr: %s", err.Error())
						}
					} else if destF.Kind() == reflect.Ptr {
						val := transFromVal(srcF.Value())
						err := destF.Set(val)
						if err != nil {
							log.GetLogger().Error(err)
							return fmt.Errorf("dest ptr: %s", err.Error())
						}
					}
					continue
				}
				err := destF.Set(srcF.Value())
				if err != nil {
					log.GetLogger().Error(err)
					return err
				}
			}
		}
	}

	return nil
}

func transFromPtr(v interface{}) interface{} {
	ty := reflect.TypeOf(v)

	var conv interface{}
	switch ty.Elem().Kind() {
	case reflect.Uint:
		conv = *v.(*uint)
	case reflect.Int:
		conv = *v.(*int)
	case reflect.Int8:
		conv = *v.(*int8)
	case reflect.Int16:
		conv = *v.(*int16)
	case reflect.Int32:
		conv = *v.(*int32)
	case reflect.Int64:
		conv = *v.(*int64)
	case reflect.Uint8:
		conv = *v.(*uint8)
	case reflect.Uint16:
		conv = *v.(*uint16)
	case reflect.Uint32:
		conv = *v.(*uint32)
	case reflect.Uint64:
		conv = *v.(*uint64)
	case reflect.Float64:
		conv = *v.(*float64)
	case reflect.Float32:
		conv = *v.(*float32)
	case reflect.String:
		conv = *v.(*string)
	case reflect.Bool:
		conv = *v.(*string)
	case reflect.Interface:
		conv = v
	}

	return conv
}

func transFromVal(v interface{}) interface{} {
	ty := reflect.TypeOf(v)

	switch ty.Kind() {
	case reflect.Uint:
		conv := v.(uint)
		return &conv
	case reflect.Int:
		conv := v.(int)
		return &conv
	case reflect.Int8:
		conv := v.(int8)
		return &conv
	case reflect.Int16:
		conv := v.(int16)
		return &conv
	case reflect.Int32:
		conv := v.(int32)
		return &conv
	case reflect.Int64:
		conv := v.(int64)
		return &conv
	case reflect.Uint8:
		conv := v.(uint8)
		return &conv
	case reflect.Uint16:
		conv := v.(uint16)
		return &conv
	case reflect.Uint32:
		conv := v.(uint32)
		return &conv
	case reflect.Uint64:
		conv := v.(uint64)
		return &conv
	case reflect.Float64:
		conv := v.(float64)
		return &conv
	case reflect.Float32:
		conv := v.(float32)
		return &conv
	case reflect.String:
		conv := v.(string)
		return &conv
	case reflect.Bool:
		conv := v.(string)
		return &conv
	case reflect.Interface:
		conv := v
		return &conv
	}
	return nil
}
