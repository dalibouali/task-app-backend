package services

import (
	"github.com/dalibouali/task-app-backend/models"
	"gorm.io/gorm"
	"fmt"

)

// GetAllUrls retrieves all Url entries from DB
func GetAllUrls(db *gorm.DB) ([]models.Url, error) {
	var urls []models.Url
	result := db.Find(&urls)
	fmt.Println("GetAllUrls called, found", len(urls), "urls")
	return urls, result.Error
}
