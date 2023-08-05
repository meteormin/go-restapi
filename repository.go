package go_restapi

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[Entity interface{}] struct {
	gormrepo.GenericRepository[Entity]
}

func NewRepository[Entity interface{}](db *gorm.DB, model Entity) *Repository[Entity] {
	return &Repository[Entity]{
		gormrepo.NewGenericRepository(db, model).Preload(clause.Associations),
	}
}

func (repo *Repository[Entity]) All(filter *Filter[Entity]) ([]Entity, error) {
	entities := make([]Entity, 0)
	db := repo.DB()
	if filter != nil {
		filter.SetEntity(repo.GenericRepository.GetModel())
		search, err := filter.Search(db)
		if err != nil {
			return entities, err
		}

		db = filter.Sort(search)
	}

	err := db.Find(&entities).Error

	return entities, err
}
