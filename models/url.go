package models
import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
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
}

