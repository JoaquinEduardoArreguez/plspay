package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                 string
	Email                string        `gorm:"unique"`
	Hashed_passwod       string        `gorm:"size:60"`
	Groups               []*Group      `gorm:"many2many:user_groups;"`
	Expenses             []Expense     `gorm:"foreignKey:OwnerID"`
	ParticipatedExpenses []*Expense    `gorm:"many2many:expense_participants;"`
	SenderTransactions   []Transaction `gorm:"foreignkey:SenderUserID"`
	ReceiverTransactions []Transaction `gorm:"foreignkey:ReceiverUserID"`
	Balances             []Balance     `gorm:"foreignKey:User"`
}

func NewUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("User name is required")
	}
	return &User{Name: name, Email: email}, nil
}
