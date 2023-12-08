package repositories

import (
	"errors"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"gorm.io/gorm"
)

var (
	InvalidCredentialsError = errors.New("invalid credentials")
	DuplicatedEmailError    = errors.New("duplicated email")
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: BaseRepository{DB: db}}
}

func (r *UserRepository) GetByName(name string, entity interface{}) *gorm.DB {
	return r.DB.Where("name = ?", name).First(entity)
}

func (r *UserRepository) GetByEmail(email string, user models.User) *gorm.DB {
	return r.DB.Where("email = ?", email).First(user)
}

func (r *UserRepository) GetUserNames() ([]string, error) {
	var users []models.User
	var userNames []string

	dbResponse := r.GetAll(&users)
	if dbResponse.Error != nil {
		return nil, dbResponse.Error
	}

	for _, user := range users {
		userNames = append(userNames, user.Name)
	}

	return userNames, nil
}

func (r *UserRepository) Authenticate(email, password string) (int, error) {
	return 0, nil
}
