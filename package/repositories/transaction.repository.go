package repositories

import "gorm.io/gorm"

type TransactionRepository struct {
	BaseRepository
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{BaseRepository: BaseRepository{DB: db}}
}
