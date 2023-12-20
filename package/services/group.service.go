package services

import (
	"errors"
	"time"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/repositories"
	"gorm.io/gorm"
)

var (
	ErrGroupNotCreated = errors.New("group not created")
)

type GroupService struct {
	groupRepository *repositories.GroupRepository
	userRepository  *repositories.UserRepository
	balanceService  *BalanceService
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		groupRepository: repositories.NewGroupRepository(db),
		userRepository:  repositories.NewUserRepository(db),
		balanceService:  NewBalanceService(db),
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
