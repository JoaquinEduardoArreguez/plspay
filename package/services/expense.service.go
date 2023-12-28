package services

import (
	"errors"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/repositories"
	"gorm.io/gorm"
)

type ExpenseService struct {
	*repositories.BaseRepository
}

func NewExpenseService(db *gorm.DB) *ExpenseService {
	return &ExpenseService{
		BaseRepository: repositories.NewBaseRepository(db),
	}
}

func (service *ExpenseService) CreateExpense(description string, amount float64, groupID uint, ownerID uint, participantsIDs []uint) (*models.Expense, error) {
	var participants []*models.User
	if err := service.DB.Find(&participants, participantsIDs).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	expense, newExpenseError := models.NewExpense(description, amount, groupID, ownerID, participants)
	if newExpenseError != nil {
		return nil, newExpenseError
	}

	if err := service.DB.Create(expense).Error; err != nil {
		return nil, err
	}

	return expense, nil
}
