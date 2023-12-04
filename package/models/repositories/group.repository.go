package repositories

import (
	"gorm.io/gorm"
)

type GroupRepository struct {
	BaseRepository
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{BaseRepository: BaseRepository{DB: db}}
}
