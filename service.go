package restapi

type Service[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] interface {
	All(filter *Filter[Entity]) ([]Res, error)
	Find(pk uint) (Res, error)
	Create(dto Req) (Res, error)
	Update(pk uint, dto Req) (Res, error)
	Patch(pk uint, dto Req) (Res, error)
	Delete(pk uint) (bool, error)
	Repo() Repository[Entity]
	Response() Res
	ServiceHook[Entity, Req, Res]
}

type GenericService[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]] struct {
	repo   Repository[Entity]
	res    Res
	events *HasServiceEvent[Entity, Req, Res]
}

func NewService[Entity interface{}, Req RequestDTO[*Entity], Res ResponseDTO[Entity]](
	repo Repository[Entity],
	resDto Res,
) Service[Entity, Req, Res] {
	return &GenericService[Entity, Req, Res]{
		repo:   repo,
		res:    resDto,
		events: &HasServiceEvent[Entity, Req, Res]{&HasMethodEvent[Entity, Req, Res]{methodEvent: nil}},
	}
}

func (s *GenericService[Entity, Req, Res]) Repo() Repository[Entity] {
	return s.repo
}

func (s *GenericService[Entity, Req, Res]) Response() Res {
	return s.res
}

func (s *GenericService[Entity, Req, Res]) All(filter *Filter[Entity]) ([]Res, error) {

	entities, err := s.repo.GetByFilter(filter)
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

	if s.features(Create).beforeCallRepo != nil {
		err = s.features(Create).beforeCallRepo.handler(dto, ent)
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

	if s.features(Create).afterCallRepo != nil {
		err = s.features(Create).afterCallRepo.handler(res, *create)
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

	if s.features(Update).beforeCallRepo != nil {
		err = s.features(Update).beforeCallRepo.handler(dto, *find)
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

	if s.features(Update).afterCallRepo != nil {
		err = s.features(Update).afterCallRepo.handler(res, *update)
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

	if s.features(Patch).beforeCallRepo != nil {
		err = s.features(Patch).beforeCallRepo.handler(dto, *find)
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

	if s.features(Patch).afterCallRepo != nil {
		err = s.features(Patch).afterCallRepo.handler(res, *update)
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

func (s *GenericService[Entity, Req, Res]) Hook() *HasServiceEvent[Entity, Req, Res] {
	return s.events
}

func (s *GenericService[Entity, Req, Res]) features(event MethodEvent) *Features[Entity, Req, Res] {
	return s.events.getMethodEvent(event)
}
