package main

import (
	"net/http"
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/dalibouali/task-app-backend/models"
	"github.com/dalibouali/task-app-backend/crawler"
	"github.com/dalibouali/task-app-backend/services"
	"github.com/dalibouali/task-app-backend/database"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// Simple health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Protect all /api routes with Basic Auth
	api := r.Group("/api", gin.BasicAuth(gin.Accounts{
		"admin":  "admin",
	}))

	// API routes
	api.GET("/urls", GetAllUrlsHandler)
	api.POST("/urls", CreateUrlHandler)
	api.PUT("/urls/:id/rerun", RerunUrlHandler)
	api.DELETE("/urls/:id", DeleteUrlHandler)

	return r
}

func GetAllUrlsHandler(c *gin.Context) {
	urls, err := services.GetAllUrls(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
// CreateUrlHandler creates a new URL entry in the database
func CreateUrlHandler(c *gin.Context) {
	var input models.Url
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Status = "queued"
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

// On Click o Rerun URL Crawl , this will set the status to "queued" again
// so that the worker can pick it up again
func RerunUrlHandler(c *gin.Context) {
	id := c.Param("id")
	var url models.Url
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	url.Status = "queued"
	database.DB.Save(&url)
	c.JSON(http.StatusOK, url)
}

func DeleteUrlHandler(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Url{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}

func main() {
	database.InitDB()
	database.DB.AutoMigrate(&models.Url{})

	// start crawler worker
	go startCrawlerWorker()

	r := setupRouter()
	r.Run(":8080")
}
// startCrawlerWorker continuously checks for queued URLs and processes them
func startCrawlerWorker() {
	for {
		var urls []models.Url
		database.DB.Where("status = ?", "queued").Find(&urls)
		for _, u := range urls {
			fmt.Println("Processing:", u.URL)
			crawler.AnalyzeUrl(database.DB, &u)
		}
		time.Sleep(5 * time.Second)
	}
}
