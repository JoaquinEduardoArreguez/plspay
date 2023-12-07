package repositories

import (
	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: BaseRepository{DB: db}}
}

func (r *UserRepository) GetByName(name string, entity interface{}) *gorm.DB {
	return r.DB.Where("name = ?", name).First(entity)
}
