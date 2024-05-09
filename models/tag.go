package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name      string
	Spendings []Spending `gorm:"many2many:spending_tag;"`
}
