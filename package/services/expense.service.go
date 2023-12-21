package services

import (
	"errors"

	utils "github.com/JoaquinEduardoArreguez/plspay/package"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/repositories"
	"gorm.io/gorm"
)

type ExpenseService struct {
	balanceService    *BalanceService
	expenseRepository *repositories.ExpenseRepository
	userRepository    *repositories.UserRepository
}

func NewExpenseService(db *gorm.DB) *ExpenseService {
	return &ExpenseService{
		balanceService:    NewBalanceService(db),
		expenseRepository: repositories.NewExpenseRepository(db),
		userRepository:    repositories.NewUserRepository(db),
	}
}

// updateUserBalances updates user balances based on the provided expense.
func (service *ExpenseService) updateUserBalances(expense models.Expense) error {
	var participantsIds []uint
	ownerIsParticipant := false

	for _, participant := range expense.Participants {
		participantsIds = append(participantsIds, participant.ID)
		if expense.OwnerID == participant.ID {
			ownerIsParticipant = true
		}
	}

	if !ownerIsParticipant {
		participantsIds = append(participantsIds, expense.OwnerID)
	}

	usersBalances, errGettingUserBalances := service.balanceService.GetBalancesByUsersAndGroup(participantsIds, expense.Group)
	if errGettingUserBalances != nil {
		return errGettingUserBalances
	}

	share := utils.RoundFloat(expense.Amount/float64(len(expense.Participants)), 2)

	var ownerBalanceIndex int
	for i, userBalance := range usersBalances {
		if userBalance.User == expense.OwnerID {
			ownerBalanceIndex = i
		} else {
			userBalance.Amount = utils.RoundFloat(userBalance.Amount-share, 2)
		}
	}

	if ownerIsParticipant {
		usersBalances[ownerBalanceIndex].Amount = utils.RoundFloat(usersBalances[ownerBalanceIndex].Amount+share*float64(len(usersBalances)-1), 2)
	} else {
		usersBalances[ownerBalanceIndex].Amount = utils.RoundFloat(usersBalances[ownerBalanceIndex].Amount+expense.Amount, 2)
	}

	if err := service.balanceService.Update(&usersBalances); err != nil {
		return err
	}

	return nil
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

	if err := service.updateUserBalances(*expense); err != nil {
		return nil, err
	}

	return expense, nil
}
