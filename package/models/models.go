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
	Date         time.Time
	Users        []*User       `gorm:"many2many:user_groups;"`
	Expenses     []Expense     `gorm:"foreignKey:Group"`
	Transactions []Transaction `gorm:"foreignKey:Group"`
	Balances     []Balance     `gorm:"foreignKey:Group"`
}

type GroupDTO struct {
	ID           uint
	CreatedAt    string
	UpdatedAt    string
	Name         string
	Date         string
	Users        string
	Expenses     []Expense
	Transactions []Transaction
}

type Balance struct {
	gorm.Model
	User   uint
	Group  uint
	Amount float64
}

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

type Expense struct {
	gorm.Model
	Description  string
	Amount       float64
	Group        uint
	OwnerID      uint
	Owner        User    `gorm:"foreignKey:OwnerID"`
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

func NewGroup(name string, owner *User, participants []*User, date time.Time) (*Group, error) {
	if owner == nil {
		return nil, errors.New("Group owner (user) is required")
	}
	if name == "" {
		return nil, errors.New("Group name is required")
	}

	groupParticipants := []*User{owner}
	groupParticipants = append(groupParticipants, participants...)

	return &Group{Name: name, Date: date, Users: groupParticipants}, nil
}

func NewUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("User name is required")
	}
	return &User{Name: name, Email: email}, nil
}

func NewExpense(description string, amount float64, groupID uint, ownerID uint, participants []*User) (*Expense, error) {
	if description == "" {
		return nil, errors.New("description is required")
	}
	if amount == 0 {
		return nil, errors.New("amount is required")
	}
	if participants == nil {
		return nil, errors.New("participants is required")
	}

	if groupID == 0 {
		return nil, errors.New("group must be defined")
	}
	if ownerID == 0 {
		return nil, errors.New("owner must be defined")
	}
	if participants == nil {
		return nil, errors.New("participants must be defined")
	}

	return &Expense{
		Description:  description,
		Amount:       amount,
		Participants: participants,
		OwnerID:      ownerID,
		Group:        groupID,
	}, nil
}

func (g *Group) ToDto() GroupDTO {
	var usernames []string
	for _, user := range g.Users {
		usernames = append(usernames, user.Name)
	}

	return GroupDTO{
		ID:           g.ID,
		CreatedAt:    g.CreatedAt.Format(time.DateTime),
		UpdatedAt:    g.UpdatedAt.Format(time.DateTime),
		Name:         g.Name,
		Date:         g.Date.Format(time.DateTime),
		Expenses:     g.Expenses,
		Transactions: g.Transactions,
		Users:        strings.Join(usernames, ","),
	}
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
