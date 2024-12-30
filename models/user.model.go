package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName string `json:"fullName" gorm:"not null"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"-" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Roll     string `json:"roll" gorm:"unique;not null"`
	Batch    int    `json:"batch" gorm:"not null"`
	FBLink   string `json:"fbLink" gorm:"not null"`
	IsAdmin  bool   `json:"isAdmin" gorm:"default:false;not null"`
}
