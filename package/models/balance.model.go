package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	User   uint
	Group  uint
	Amount float64
}
