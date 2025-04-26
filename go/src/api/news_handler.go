package api

import (
	"net/http"

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
	result := h.db.Preload("Reactions").Find(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	newsResponses := models.ToGetAllNewsResponse(news)
	c.JSON(http.StatusOK, newsResponses)
}

func (h *NewsHandler) GetNewsReactionsById(c *gin.Context) {
	newsId := c.Query("id")

	var reactions []models.NewsReaction
	result := h.db.Where("news_id = ?", newsId).Find(&reactions)
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
