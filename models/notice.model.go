package models

import "gorm.io/gorm"

type Notice struct {
	gorm.Model `json:"-"`
	Title      string `json:"title" gorm:"not null"`                      // Title of the notice
	Date       string `json:"date" gorm:"not null"`                       // Date of the notice
	Content    string `json:"content" gorm:"not null"`                    // Content of the notice
	Files      []File `json:"files,omitempty" gorm:"foreignKey:NoticeID"` // Associated files
}
