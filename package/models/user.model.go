package models

import (
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
	IsGuest              bool          `gorm:"not null;default:false"`
}
