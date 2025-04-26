package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping"
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

	now := time.Now()
	result := h.db.Where("DATE(created_at) >= DATE(?)", now).Limit(100).Preload("Reactions").Find(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if len(news) == 0 {
		yesterday := now.AddDate(0, 0, -1)
		result = h.db.Where("DATE(created_at) >= DATE(?)", yesterday).Limit(100).Preload("Reactions").Find(&news)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	rand.Shuffle(len(news), func(i, j int) {
		news[i], news[j] = news[j], news[i]
	})

	newsResponses := models.ToGetAllNewsResponse(news)
	c.JSON(http.StatusOK, newsResponses)
}

func (h *NewsHandler) GetNewsReactionsById(c *gin.Context) {
	newsId := c.Query("id")

	var reactions []models.NewsReaction
	result := h.db.Where("news_id = ?", newsId).Limit(100).Find(&reactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, reactions)
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
			reactions, err := siteScraper.ScrapeReactions(news[i])
			if err != nil {
				println("Error scraping reactions:", err.Error())
			}
			news[i].Reactions = reactions
			h.db.Create(&news[i])
		}

		println("News scraping completed successfully, articles saved:", len(news))
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News scraping has been started in the background",
	})
}
