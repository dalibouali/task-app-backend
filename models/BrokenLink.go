package models

import "gorm.io/gorm"

type BrokenLink struct {
	gorm.Model
	URL        string
	StatusCode int
	UrlID      uint 
}
