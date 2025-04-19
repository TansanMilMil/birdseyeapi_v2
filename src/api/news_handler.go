package api

import (
	"net/http"
	"strconv"

	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/src/scraping"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewsHandler handles news-related API requests
type NewsHandler struct {
	db *gorm.DB
}

// NewNewsHandler creates a new NewsHandler
func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{db: db}
}

// GetAllNews returns all news articles
func (h *NewsHandler) GetAllNews(c *gin.Context) {
	var news []models.News
	result := h.db.Preload("Reactions").Find(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

// GetNewsById returns a news article by ID
func (h *NewsHandler) GetNewsById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var news models.News
	result := h.db.Preload("Reactions").First(&news, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News article not found"})
		return
	}

	c.JSON(http.StatusOK, news)
}

// CreateNews creates a new news article
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Create(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, news)
}

// ScrapeNews triggers scraping of news articles
func (h *NewsHandler) ScrapeNews(c *gin.Context) {
	// Create a new site scraper
	siteScraper := scraping.NewSiteScraping()
	
	// Start the scraping process in a goroutine to avoid timeout
	go func() {
		// Scrape news articles in the background
		news, err := siteScraper.ScrapeNews()
		if err != nil {
			// Just log the error since we can't return it to the client anymore
			// In a production app, you might want to use a proper logging system
			// or error tracking service here
			println("Error scraping news:", err.Error())
			return
		}

		// Save news articles to the database
		for i := range news {
			h.db.Create(&news[i])
		}
		
		println("News scraping completed successfully, articles saved:", len(news))
	}()

	// Return an immediate response to the client
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News scraping has been started in the background",
	})
}

// SummarizeNews summarizes a news article
func (h *NewsHandler) SummarizeNews(c *gin.Context) {
	var req models.NewsSummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find news article by URL
	var news models.News
	result := h.db.Where("article_url = ?", req.URL).First(&news)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News article not found"})
		return
	}

	// Use the AI service to summarize the text
	summarizer := scraping.NewOpenAISummarizer()
	summarizedText, err := summarizer.Summarize(news.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the news article with the summarized text
	news.SummarizedText = summarizedText
	h.db.Save(&news)

	c.JSON(http.StatusOK, news)
}