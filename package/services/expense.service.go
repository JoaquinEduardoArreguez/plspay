package services

import (
	"gorm.io/gorm"
)

type ExpenseService struct {
	balanceService *BalanceService
}

func NewExpenseService(database *gorm.DB) *ExpenseService {
	return &ExpenseService{balanceService: NewBalanceService(database)}
}
