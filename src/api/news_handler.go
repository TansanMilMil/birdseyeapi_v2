package api

import (
	"net/http"
	"strconv"

	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/src/scraping"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsHandler struct {
	db *gorm.DB
}

func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{db: db}
}

func (h *NewsHandler) GetAllNews(c *gin.Context) {
	var news []models.News
	result := h.db.Preload("Reactions").Find(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

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

func (h *NewsHandler) ScrapeNews(c *gin.Context) {
	siteScraper := scraping.NewSiteScraping()
	
	go func() {
		news, err := siteScraper.ScrapeNews()
		if err != nil {
			println("Error scraping news:", err.Error())
			return
		}

		for i := range news {
			h.db.Create(&news[i])
		}
		
		println("News scraping completed successfully, articles saved:", len(news))
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News scraping has been started in the background",
	})
}

func (h *NewsHandler) SummarizeNews(c *gin.Context) {
	var req models.NewsSummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var news models.News
	result := h.db.Where("article_url = ?", req.URL).First(&news)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News article not found"})
		return
	}

	summarizer := scraping.NewOpenAISummarizer()
	summarizedText, err := summarizer.Summarize(news.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	news.SummarizedText = summarizedText
	h.db.Save(&news)

	c.JSON(http.StatusOK, news)
}