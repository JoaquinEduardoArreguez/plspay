package repositories

import (
	"errors"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	InvalidCredentialsError = errors.New("invalid credentials")
	//DuplicatedEmailError    = errors.New("duplicated email")
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

func (r *UserRepository) FindByNames(names []string, users *[]models.User) *gorm.DB {
	return r.DB.Where("name IN ?", names).Find(users)
}

func (r *UserRepository) GetByEmail(email string, user *models.User) *gorm.DB {
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

func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	dbResponse := r.GetAll(&users)
	if dbResponse.Error != nil {
		return nil, dbResponse.Error
	}
	return users, nil
}

func (r *UserRepository) Authenticate(email, password string) (int, error) {
	var user models.User

	dbResponse := r.GetByEmail(email, &user)
	if errors.Is(dbResponse.Error, gorm.ErrRecordNotFound) {
		return 0, InvalidCredentialsError
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hashed_passwod), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, InvalidCredentialsError
	} else if err != nil {
		return 0, err
	}

	return int(user.ID), nil
}

func (r *UserRepository) CreateUser(name, email, password string) (*gorm.DB, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{Name: name, Email: email, Hashed_passwod: string(hashedPassword)}

	dbResponse := r.Create(newUser)
	if dbResponse.Error != nil {
		return nil, dbResponse.Error
	}

	return dbResponse, nil
}
