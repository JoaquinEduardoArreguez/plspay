package services

import (
	"errors"
	"math"
	"time"

	utils "github.com/JoaquinEduardoArreguez/plspay/package"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/repositories"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrGroupNotCreated = errors.New("group not created")
	ErrGroupNotFound   = errors.New("group not found")
)

type GroupService struct {
	groupRepository       *repositories.GroupRepository
	userRepository        *repositories.UserRepository
	transactionRepository *repositories.TransactionRepository
	balanceService        *BalanceService
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		groupRepository:       repositories.NewGroupRepository(db),
		userRepository:        repositories.NewUserRepository(db),
		transactionRepository: repositories.NewTransactionRepository(db),
		balanceService:        NewBalanceService(db),
	}
}

func (service *GroupService) CreateGroup(name string, owner *models.User, participantNames []string, date time.Time) (*models.Group, error) {
	var participants []*models.User
	if err := service.userRepository.FindByNames(participantNames, &participants); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	group, newGroupError := models.NewGroup(name, owner, participants, date)
	if newGroupError != nil {
		return nil, newGroupError
	}

	if err := service.groupRepository.Create(&group); err != nil {
		return nil, ErrGroupNotCreated
	}

	userBalances := make([]*models.Balance, len(participants))
	for i, participant := range participants {
		userBalances[i] = &models.Balance{User: participant.ID, Group: group.ID, Amount: 0}
	}

	if err := service.balanceService.Create(&userBalances); err != nil {
		return nil, err
	}

	return group, nil
}

func (service *GroupService) CreateTransactions(groupId uint) ([]*models.Transaction, error) {
	if err := service.clearTransactions(groupId); err != nil {
		return nil, err
	}

	group := &models.Group{}
	if err := service.groupRepository.GetByID(groupId, &group, "Balances"); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGroupNotFound
	} else if err != nil {
		return nil, err
	}

	var transactions []*models.Transaction
	transactions = service.generateTransactions(group.Balances, transactions)

	if len(transactions) > 0 {
		if err := service.transactionRepository.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"amount"}),
		}).Create(&transactions).Error; err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func (service *GroupService) generateTransactions(remainingBalances []*models.Balance, transactions []*models.Transaction) []*models.Transaction {
	if !service.balancesSettled(remainingBalances) {
		var sender, receiver *models.Balance
		var minBalance, maxBalance float64

		// Find sender and receiver
		for _, remainingBalance := range remainingBalances {
			if remainingBalance.Amount < minBalance {
				sender = remainingBalance
			} else if remainingBalance.Amount > maxBalance {
				receiver = remainingBalance
			}
		}

		// Calculate transaction amount
		var transactionAmount float64 = math.Min(math.Abs(sender.Amount), math.Abs(receiver.Amount))

		// Add transaction
		transactions = append(transactions, &models.Transaction{
			Amount:         transactionAmount,
			SenderUserID:   sender.User,
			ReceiverUserID: receiver.User,
			Group:          remainingBalances[0].Group,
		})

		// Update user balances
		sender.Amount = math.Copysign(utils.RoundFloat(math.Abs(sender.Amount)-transactionAmount, 2), sender.Amount)
		receiver.Amount = math.Copysign(utils.RoundFloat(math.Abs(receiver.Amount)-transactionAmount, 2), receiver.Amount)

		// Calculate next transaction
		return service.generateTransactions(remainingBalances, transactions)
	}

	return transactions
}

func (service *GroupService) balancesSettled(balances []*models.Balance) bool {
	for _, balance := range balances {
		if math.Abs(balance.Amount) > 1 {
			return false
		}
	}
	return true
}

func (service *GroupService) clearTransactions(groupID uint) error {
	return service.transactionRepository.DB.Where(models.Transaction{Group: groupID}).Delete(&models.Transaction{}).Error
}

func (service *GroupService) DeleteGroup(groupID uint) error {
	return service.groupRepository.Delete(groupID, &models.Group{})
}
