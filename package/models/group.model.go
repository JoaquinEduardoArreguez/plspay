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
	Expenses     []Expense     `gorm:"foreignKey:GroupID"`
	Transactions []Transaction `gorm:"foreignKey:Group"`
	Balances     []*Balance    `gorm:"foreignKey:Group"`
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
