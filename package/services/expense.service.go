package services

import (
	"errors"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/repositories"
	"gorm.io/gorm"
)

type ExpenseService struct {
	*repositories.BaseRepository
	balanceService    *BalanceService
	expenseRepository *repositories.ExpenseRepository
	userRepository    *repositories.UserRepository
}

func NewExpenseService(db *gorm.DB) *ExpenseService {
	return &ExpenseService{
		BaseRepository:    repositories.NewBaseRepository(db),
		balanceService:    NewBalanceService(db),
		expenseRepository: repositories.NewExpenseRepository(db),
		userRepository:    repositories.NewUserRepository(db),
	}
}

func (service *ExpenseService) CreateExpense(description string, amount float64, groupID uint, ownerID uint, participantsIDs []uint) (*models.Expense, error) {
	var participants []*models.User
	if err := service.userRepository.DB.Find(&participants, participantsIDs).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	expense, newExpenseError := models.NewExpense(description, amount, groupID, ownerID, participants)
	if newExpenseError != nil {
		return nil, newExpenseError
	}

	if err := service.expenseRepository.Create(expense); err != nil {
		return nil, err
	}

	return expense, nil
}
