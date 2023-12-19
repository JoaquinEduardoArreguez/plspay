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

// GetBalanceByUserAndGroup retrieves the balance for a user within a specific group.
func (s *BalanceService) GetBalanceByUserAndGroup(userID, groupID uint) (*models.Balance, error) {
	var balance models.Balance
	if err := s.DB.Where("user = ? AND group = ?", userID, groupID).First(&balance).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBalanceNotFound
	} else if err != nil {
		return nil, err
	}

	return &balance, nil
}

// UpdateBalance updates the balance for a user within a specific group.
func (s *BalanceService) UpdateBalance(currentBalance *models.Balance, newAmount float64) error {
	currentBalance.Amount = newAmount

	if err := s.DB.Save(currentBalance).Error; err != nil {
		return fmt.Errorf("%w: %v", ErrBalanceNotSaved, err)
	}

	return nil
}
