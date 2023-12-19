package repositories

import "gorm.io/gorm"

type ExpenseRepository struct {
	BaseRepository
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{BaseRepository: BaseRepository{DB: db}}
}
