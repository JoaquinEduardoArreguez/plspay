package models

import (
	"errors"
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
