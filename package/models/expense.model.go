package models

import (
	"errors"

	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Description  string
	Amount       float64
	Group        uint
	OwnerID      uint
	Owner        User    `gorm:"foreignKey:OwnerID"`
	Participants []*User `gorm:"many2many:expense_participants;"`
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
