package services

import (
	"errors"
	"fmt"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/repositories"
	"gorm.io/gorm"
)

var (
	ErrBalanceNotFound = errors.New("balance not found")
	ErrBalanceNotSaved = errors.New("balance not saved")
)

type BalanceService struct {
	*repositories.BaseRepository
}

func NewBalanceService(db *gorm.DB) *BalanceService {
	return &BalanceService{BaseRepository: repositories.NewBaseRepository(db)}
}

// UpdateBalance updates the balance for a user within a specific group.
func (service *BalanceService) UpdateBalance(currentBalance *models.Balance, newAmount float64) error {
	currentBalance.Amount = newAmount

	if err := service.DB.Save(currentBalance).Error; err != nil {
		return fmt.Errorf("%w: %v", ErrBalanceNotSaved, err)
	}

	return nil
}

// GetBalancesByUsersAndGroup retrieves balances for multiple users within a specific group.
func (service *BalanceService) GetBalancesByUsersAndGroup(userIDs []uint, groupID uint) ([]*models.Balance, error) {
	var balances []*models.Balance

	queryConditions := map[string]interface{}{
		"user":  userIDs,
		"group": groupID,
	}

	if err := service.DB.Where(queryConditions).Find(&balances).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBalanceNotFound, err)
	}
	return balances, nil
}
