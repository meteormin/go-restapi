package restapi

import "github.com/miniyus/gofiber"

type Service[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	All(filter *Filter[Entity]) ([]Res, error)
	Find(pk uint) (Res, error)
	Create(dto Req) (Res, error)
	Update(pk uint, dto Req) (Res, error)
	Patch(pk uint, dto Req) (Res, error)
	Delete(pk uint) (bool, error)
	Repo() *Repository[Entity]
	Response() Res
}

type GenericService[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	repo   *Repository[Entity]
	res    Res
	events HasMethodEvent
}

type GenericRepository[Entity interface{}] struct {
	Repository[Entity]
}

func NewService[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](
	repo *Repository[Entity],
	resDto Res,
) Service[Entity, Req, Res] {
	return &GenericService[Entity, Req, Res]{
		repo:   repo,
		res:    resDto,
		events: HasMethodEvent{},
	}
}

func (s *GenericService[Entity, Req, Res]) Repo() *Repository[Entity] {
	return s.repo
}

func (s *GenericService[Entity, Req, Res]) Response() Res {
	return s.res
}

func (s *GenericService[Entity, Req, Res]) All(filter *Filter[Entity]) ([]Res, error) {
	entities, err := s.repo.All(filter)
	if err != nil {
		return nil, err
	}

	res := make([]Res, 0)
	for _, ent := range entities {
		temp := DeepCopy(s.res).(Res)
		err = temp.FromEntity(ent)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (s *GenericService[Entity, Req, Res]) Find(pk uint) (Res, error) {
	res := s.res
	entity, err := s.repo.Find(pk)
	//entity, err := s.repo.Find(pk)
	if err != nil {
		return res, err
	}

	err = res.FromEntity(*entity)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *GenericService[Entity, Req, Res]) Create(dto Req) (Res, error) {
	res := s.res
	ent := s.repo.GetModel()
	err := dto.ToEntity(&ent)
	if err != nil {
		return res, err
	}

	if s.features(Create).BeforeCallRepo != nil {
		err = s.features(Create).BeforeCallRepo.handler(dto, ent)
		if err != nil {
			return res, err
		}
	}

	create, err := s.repo.Create(ent)
	if err != nil {
		return res, err
	}

	err = res.FromEntity(*create)
	if err != nil {
		return res, err
	}

	if s.features(Create).AfterCallRepo != nil {
		err = s.features(Create).AfterCallRepo.handler(res, create)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (s *GenericService[Entity, Req, Res]) Update(pk uint, dto Req) (Res, error) {
	res := s.res
	find, err := s.repo.Find(pk)
	if err != nil {
		return res, err
	}

	err = dto.ToEntity(find)
	if err != nil {
		return res, err
	}

	if s.features(Update).BeforeCallRepo != nil {
		err = s.features(Update).BeforeCallRepo.handler(dto, find)
		if err != nil {
			return res, err
		}
	}

	update, err := s.repo.Update(pk, *find)
	if err != nil {
		return res, err
	}

	err = res.FromEntity(*update)
	if err != nil {
		return res, err
	}

	if s.features(Update).AfterCallRepo != nil {
		err = s.features(Update).AfterCallRepo.handler(res, update)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (s *GenericService[Entity, Req, Res]) Patch(pk uint, dto Req) (Res, error) {
	res := s.res
	find, err := s.repo.Find(pk)
	if err != nil {
		return res, err
	}

	err = dto.ToEntity(find)
	if err != nil {
		return res, err
	}

	if s.features(Patch).BeforeCallRepo != nil {
		err = s.features(Patch).BeforeCallRepo.handler(dto, find)
		if err != nil {
			return res, err
		}
	}

	update, err := s.repo.Update(pk, *find)
	if err != nil {
		return res, err
	}

	err = res.FromEntity(*update)
	if err != nil {
		return res, err
	}

	if s.features(Patch).AfterCallRepo != nil {
		err = s.features(Patch).AfterCallRepo.handler(res, update)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (s *GenericService[Entity, Req, Res]) Delete(pk uint) (bool, error) {
	b, err := s.repo.Delete(pk)
	if err != nil {
		return false, err
	}
	return b, nil
}

func (s *GenericService[Entity, Req, Res]) Hook() HasMethodEvent {
	return s.events
}

func (s *GenericService[Entity, Req, Res]) features(event MethodEvent) *Features {
	feat, err := s.events.getMethodEvent(event)
	gofiber.Log().Error(err)

	return feat
}
