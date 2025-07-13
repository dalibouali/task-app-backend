package models

import (
	"gorm.io/gorm"
	"time"
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
	BrokenLinksCount int `json:"brokenLinksCount"`
	HasLoginForm  bool   `json:"hasLoginForm"`
	Status        string `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"`
	BrokenLinksList []BrokenLink `gorm:"foreignKey:UrlID;references:ID" json:"brokenLinksList"`
}
