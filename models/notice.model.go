package models

import "gorm.io/gorm"

type Notice struct {
	gorm.Model
	Title   string `json:"title"`                                      // Title of the notice
	Date    string `json:"date"`                                       // Date of the notice
	Content string `json:"content"`                                    // Content of the notice
	Files   []File `json:"files,omitempty" gorm:"foreignKey:NoticeID"` // Associated files
}
