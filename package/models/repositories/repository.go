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

func (r *BaseRepository) GetByID(id uint, entity interface{}) *gorm.DB {
	return r.DB.First(entity, id)
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
