package models

import (
	"gorm.io/gorm"
)

type Spending struct {
	gorm.Model
	ChatId      int
	MessageId   int
	Cost        float64
	Description string
	Tags        []Tag `gorm:"many2many:spending_tag;"`
}
