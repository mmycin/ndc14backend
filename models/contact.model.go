package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	Name    string `json:"name"`
	Email   string `gorm:"unique" json:"email"`
	Roll    string `json:"roll"`
	Message string `json:"message"`
}
