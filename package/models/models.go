package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name         string
	Description  string
	Users        []*User       `gorm:"many2many:user_groups;"`
	Expenses     []Expense     `gorm:"foreignKey:Group"`
	Transactions []Transaction `gorm:"foreignKey:Group"`
}

type GroupDTO struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Description  string
	Users        string
	Expenses     []Expense
	Transactions []Transaction
}

type User struct {
	gorm.Model
	Name                 string
	Email                string
	Groups               []*Group      `gorm:"many2many:user_groups;"`
	Expenses             []Expense     `gorm:"foreignKey:Owner"`
	ParticipatedExpenses []*Expense    `gorm:"many2many:expense_participants;"`
	SenderTransactions   []Transaction `gorm:"foreignkey:SenderUserID"`
	ReceiverTransactions []Transaction `gorm:"foreignkey:ReceiverUserID"`
}

type Expense struct {
	gorm.Model
	Description  string
	Amount       float64
	Group        uint
	Owner        uint
	Participants []*User `gorm:"many2many:expense_participants;"`
}

type Transaction struct {
	gorm.Model
	Amount         float64
	Group          uint
	SenderUserID   uint
	ReceiverUserID uint
	Sender         User `gorm:"foreignkey:SenderUserID"`
	Receiver       User `gorm:"foreignkey:ReceiverUserID"`
}

func NewGroup(owner *User, name string, description string) (*Group, error) {
	if owner == nil {
		return nil, errors.New("Group owner (user) is required")
	}
	if name == "" {
		return nil, errors.New("Group name is required")
	}

	return &Group{Name: name, Description: description, Users: []*User{owner}}, nil
}

func NewUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("User name is required")
	}
	return &User{Name: name, Email: email}, nil
}

func (g *Group) ToDto() GroupDTO {
	var usernames []string
	for _, user := range g.Users {
		usernames = append(usernames, user.Name)
	}

	return GroupDTO{
		CreatedAt:    g.CreatedAt,
		UpdatedAt:    g.UpdatedAt,
		Name:         g.Name,
		Description:  g.Description,
		Expenses:     g.Expenses,
		Transactions: g.Transactions,
		Users:        strings.Join(usernames, ","),
	}
}
