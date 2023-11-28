package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name         string
	Description  string
	Users        []*User       `gorm:"many2many:user_groups;"`
	Expenses     []Expense     `gorm:"foreignKey:Group"`
	Transactions []Transaction `gorm:"foreignKey:Group"`
}

type User struct {
	gorm.Model
	Name         string
	Email        string
	Groups       []*Group      `gorm:"many2many:user_groups;"`
	Expenses     []Expense     `gorm:"foreignKey:Owner"`
	Transactions []Transaction `gorm:"foreignKey:Sender"`
	Participant  uint
}

type Expense struct {
	gorm.Model
	Description  string
	Amount       float64
	Group        uint
	Owner        uint
	Participants []User `gorm:"foreignKey:Participant"`
}

type Transaction struct {
	gorm.Model
	Amount   float64
	Group    uint
	Sender   uint
	Receiver uint
	User     User `gorm:"foreignKey:Receiver"`
}
