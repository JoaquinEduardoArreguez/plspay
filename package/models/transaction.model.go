package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount         float64
	Group          uint
	SenderUserID   uint
	ReceiverUserID uint
	Sender         User `gorm:"foreignkey:SenderUserID"`
	Receiver       User `gorm:"foreignkey:ReceiverUserID"`
}
