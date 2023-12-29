package services

import (
	"errors"
	"math"
	"time"

	utils "github.com/JoaquinEduardoArreguez/plspay/package"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrGroupNotCreated = errors.New("group not created")
	ErrGroupNotFound   = errors.New("group not found")
)

type GroupService struct {
	DB             *gorm.DB
	balanceService *BalanceService
}

func NewGroupService(database *gorm.DB) *GroupService {
	return &GroupService{
		DB:             database,
		balanceService: NewBalanceService(database),
	}
}

func (service *GroupService) GetGroupById(dest *models.Group, groupId int, relations ...string) error {
	query := service.DB.Where("id = ?", groupId)

	for _, relation := range relations {
		query.Preload(relation)
	}

	return query.First(dest).Error
}

func (service *GroupService) CreateGroup(name string, owner *models.User, participantNames []string, date time.Time) (*models.Group, error) {
	var participants []*models.User

	if err := service.DB.Where("name IN ?", participantNames).Find(&participants).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	group, newGroupError := models.NewGroup(name, owner, participants, date)
	if newGroupError != nil {
		return nil, newGroupError
	}

	if err := service.DB.Create(&group).Error; err != nil {
		return nil, ErrGroupNotCreated
	}

	return group, nil
}

func (service *GroupService) CreateTransactions(groupId uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	//Get group by ID
	var group models.Group
	if err := service.DB.Preload("Users").Preload("Expenses").Preload("Expenses.Participants").First(&group, groupId).Error; err != nil {
		return nil, err
	}

	if len(group.Expenses) == 0 {
		return transactions, nil
	}

	usersBalances := make([]*models.Balance, len(group.Users))
	for i, participant := range group.Users {
		usersBalances[i] = &models.Balance{User: participant.ID, Group: group.ID, Amount: 0}
	}

	dbTransaction := service.DB.Begin()

	// Delete prev transactions, if applicable
	if err := dbTransaction.Where("\"group\" = ?", groupId).Delete(&models.Transaction{}).Error; err != nil {
		dbTransaction.Rollback()
		return nil, err
	}

	// Delete prev balances, if applicable
	if err := dbTransaction.Where("\"group\" = ?", groupId).Delete(&models.Balance{}).Error; err != nil {
		dbTransaction.Rollback()
		return nil, err
	}

	if err := dbTransaction.Create(&usersBalances).Error; err != nil {
		dbTransaction.Rollback()
		return nil, err
	}

	if err := dbTransaction.Commit().Error; err != nil {
		dbTransaction.Rollback()
		return nil, err
	}

	for _, expense := range group.Expenses {
		if err := service.updateUserBalances(expense); err != nil {
			return nil, err
		}
	}

	if err := service.DB.Preload("Balances").First(&group, groupId).Error; err != nil {
		return nil, err
	}

	transactions = service.generateTransactions(group.Balances, transactions)

	if len(transactions) > 0 {
		if err := service.DB.Clauses(clause.OnConflict{
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

func (service *GroupService) DeleteGroup(groupID uint) error {
	dbTransaction := service.DB.Begin()

	if err := dbTransaction.Where("\"group\" = ?", groupID).Delete(&models.Transaction{}).Error; err != nil {
		dbTransaction.Rollback()
	}

	if err := dbTransaction.Where("\"group\" = ?", groupID).Delete(&models.Balance{}).Error; err != nil {
		dbTransaction.Rollback()
	}

	if err := dbTransaction.Where("\"group_id\" = ?", groupID).Delete(&models.Expense{}).Error; err != nil {
		dbTransaction.Rollback()
	}

	if err := dbTransaction.Delete(&models.Group{}, groupID).Error; err != nil {
		dbTransaction.Rollback()
	}

	return dbTransaction.Commit().Error
}

func (service *GroupService) updateUserBalances(expense models.Expense) error {
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

	usersBalances, errGettingUserBalances := service.balanceService.GetBalancesByUsersAndGroup(participantsIds, expense.GroupID)
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

	if err := service.DB.Save(&usersBalances).Error; err != nil {
		return err
	}

	return nil
}
