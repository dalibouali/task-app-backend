package main

import (
	"github.com/dalibouali/task-app-backend/models"
)

func SeedDB() {
	// Check if DB already has data
	var count int64
	DB.Model(&models.Url{}).Count(&count)
	if count > 0 {
		return 
	}

	// Insert dummy URLs
	DB.Create(&models.Url{
		URL:           "https://example.com",
		HtmlVersion:   "HTML5",
		Title:         "Example Domain",
		H1Count:       2,
		H2Count:       3,
		InternalLinks: 10,
		ExternalLinks: 5,
		BrokenLinksCount:   1,
		HasLoginForm:  false,
		Status:        "done",
	})

	DB.Create(&models.Url{
		URL:           "https://test.com",
		HtmlVersion:   "HTML4.01",
		Title:         "Test Site",
		H1Count:       1,
		H2Count:       4,
		InternalLinks: 7,
		ExternalLinks: 8,
		BrokenLinksCount:   0,
		HasLoginForm:  true,
		Status:        "done",
	})
}
