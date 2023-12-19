package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                 string
	Email                string        `gorm:"unique"`
	Hashed_passwod       string        `gorm:"size:60"`
	Groups               []*Group      `gorm:"many2many:user_groups;"`
	Expenses             []Expense     `gorm:"foreignKey:Owner"`
	ParticipatedExpenses []*Expense    `gorm:"many2many:expense_participants;"`
	SenderTransactions   []Transaction `gorm:"foreignkey:SenderUserID"`
	ReceiverTransactions []Transaction `gorm:"foreignkey:ReceiverUserID"`
	Balances             []Balance     `gorm:"foreignKey:User"`
}

type UserDTO struct {
	ID        uint
	CreatedAt string
	UpdatedAt string
	Name      string
	Email     string
}

func NewUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("User name is required")
	}
	return &User{Name: name, Email: email}, nil
}

func (u *User) ToDto() UserDTO {
	return UserDTO{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Format(time.DateTime),
		UpdatedAt: u.UpdatedAt.Format(time.DateTime),
		Name:      u.Name,
		Email:     u.Email,
	}
}
