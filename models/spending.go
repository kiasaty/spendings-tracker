package models

import (
	"time"

	"gorm.io/gorm"
)

type Spending struct {
	gorm.Model
	ChatId      int64
	MessageId   int
	Cost        float64
	Description string
	SpentAt     time.Time
	Tags        []Tag `gorm:"many2many:spending_tag;"`
}
