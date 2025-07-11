package models

import (
	"gorm.io/gorm"
)

type Url struct {
	ID            uint   `json:"id"`
	URL           string `json:"url"`
	HtmlVersion   string `json:"htmlVersion"`
	Title         string `json:"title"`
	H1Count       int    `json:"h1Count"`
	H2Count       int    `json:"h2Count"`
	InternalLinks int    `json:"internalLinks"`
	ExternalLinks int    `json:"externalLinks"`
	BrokenLinks   int    `json:"brokenLinks"`
	HasLoginForm  bool   `json:"hasLoginForm"`
	Status        string `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"`
}
