package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/aws"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
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
	news = h.getNews(now, c)

	if len(news) == 0 {
		yesterday := now.AddDate(0, 0, -1)
		news = h.getNews(yesterday, c)
	}

	h.shuffle(news)

	newsResponses := models.ToGetAllNewsResponse(news)

	c.JSON(http.StatusOK, newsResponses)
}

func (h *NewsHandler) getNews(t time.Time, c *gin.Context) []models.News {
	var news []models.News

	result := h.db.
		Where("DATE(created_at) >= DATE(?)", t).
		Limit(100).
		Preload("Reactions").
		Find(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return nil
	}
	return news
}

func (h *NewsHandler) shuffle(news []models.News) {
	rand.Shuffle(len(news), func(i, j int) {
		news[i], news[j] = news[j], news[i]
	})
}

func (h *NewsHandler) GetNewsReactionsById(c *gin.Context) {
	newsId := c.Param("news-id")

	var reactions []models.NewsReaction
	result := h.db.
		Where("news_id = ?", newsId).
		Limit(100).
		Find(&reactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, reactions)
}

func (h *NewsHandler) Scrape(c *gin.Context) {
	siteScraper := scraping.NewSiteScraping()

	go func() {
		news := h.scrapeNews(siteScraper, c)

		h.scrapeReactions(news, siteScraper)

		h.invalidateApiCache()
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News scraping has been started in the background",
	})
}

func (h *NewsHandler) scrapeNews(siteScraper *scraping.SiteScraping, c *gin.Context) []models.News {
	news, err := siteScraper.ScrapeNews()
	if err != nil {
		println("Error scraping news:", err.Error())
		return nil
	}

	return news
}

func (h *NewsHandler) scrapeReactions(news []models.News, siteScraper *scraping.SiteScraping) []models.News {
	for i := range news {
		reactions, err := siteScraper.ScrapeReactions(news[i])
		if err != nil {
			println("Error scraping reactions:", err.Error())
		}
		news[i].Reactions = reactions
		h.db.Create(&news[i])
	}

	println("News scraping completed successfully, articles saved:", len(news))

	return news
}

func (h *NewsHandler) invalidateApiCache() {
	err := aws.CreateInvalidation(
		env.GetEnv("AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID", ""),
		[]string{
			"/news/today-news",
			"/news/news-reactions/*",
		})
	if err != nil {
		println("Error creating CloudFront invalidation:", err.Error())
	} else {
		println("CloudFront invalidation created successfully")
	}
}
