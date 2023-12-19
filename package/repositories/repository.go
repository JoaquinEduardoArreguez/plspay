package repositories

import "gorm.io/gorm"

type Repository interface {
	GetByID(id uint, entity interface{}) error
	Create(entity interface{}) error
	Update(entity interface{}) error
	Delete(id uint, entity interface{}) error
}

type BaseRepository struct {
	DB *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

func (r *BaseRepository) GetAll(entities interface{}) *gorm.DB {
	return r.DB.Find(entities)
}

func (r *BaseRepository) GetByID(id uint, entity interface{}, relations ...string) *gorm.DB {
	query := r.DB.Model(entity).Where("id = ?", id)

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	return query.First(entity)
}

func (r *BaseRepository) Create(entity interface{}) *gorm.DB {
	return r.DB.Create(entity)
}

func (r *BaseRepository) Update(entity interface{}) *gorm.DB {
	return r.DB.Save(entity)
}

func (r *BaseRepository) Delete(id uint, entity interface{}) *gorm.DB {
	return r.DB.Delete(entity, id)
}
