package services

import (
	"errors"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{DB: database}
}

func (service *UserService) CreateUser(name, email, password string) (*models.User, error) {
	hashedPassword, hashedPasswordError := bcrypt.GenerateFromPassword([]byte(password), 12)
	if hashedPasswordError != nil {
		return nil, hashedPasswordError
	}

	user := &models.User{Name: name, Email: email, Hashed_passwod: string(hashedPassword)}

	if err := service.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) Authenticate(email, password string) (int, error) {
	var user models.User

	if err := service.DB.Where("email = ?", email).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hashed_passwod), []byte(password)); err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return int(user.ID), nil
}

func (service *UserService) GetUserById(dest *models.User, userId int, relations ...string) error {
	query := service.DB.Where("id = ?", userId)

	for _, relation := range relations {
		query.Preload(relation)
	}

	return query.First(dest).Error
}

func (service *UserService) GetUserSuggestionsByEmail(input string) (map[string]string, error) {
	var matchingUsers []models.User
	if err := service.DB.Where("email LIKE ?", "%"+input+"%").Find(&matchingUsers).Error; err != nil {
		return nil, err
	}

	suggestions := make(map[string]string, len(matchingUsers))
	for _, matchingUser := range matchingUsers {
		suggestions[matchingUser.Email] = matchingUser.Name
	}

	return suggestions, nil
}
