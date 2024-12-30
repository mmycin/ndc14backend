package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Filename string `json:"filename"`  // Name of the file
	Index    string `json:"index"`     // Index for ordering or identification
	Format   Format `json:"format"`    // File format
	NoticeID uint   `json:"notice_id"` // Foreign key linking to Notice
}

// Format is an enum-like custom type for file formats.
type Format string

// Constants for supported file formats.
const (
	PDF  Format = "PDF"
	PNG  Format = "PNG"
	JPG  Format = "JPG"
	JPEG Format = "JPEG"
	WEBP Format = "WEBP"
)
