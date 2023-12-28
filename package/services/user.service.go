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
