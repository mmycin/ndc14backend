package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email" gorm:"unique"`
	Roll     string `json:"roll" gorm:"unique"`
	Batch    int    `json:"batch"`
	FBLink   string `json:"fbLink"`
	IsAdmin  bool   `json:"isAdmin" gorm:"default:false"`
}
