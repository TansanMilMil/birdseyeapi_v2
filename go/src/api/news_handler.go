package api

import (
	"log"
	"net/http"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/cache"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	repo "github.com/birdseyeapi/birdseyeapi_v2/go/src/repository"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/util/slice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsHandler struct {
	db           *gorm.DB
	newsRepo     *repo.NewsRepository
	reactionRepo *repo.NewsReactionRepository
}

func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{
		db:           db,
		newsRepo:     repo.NewNewsRepository(db),
		reactionRepo: repo.NewNewsReactionRepository(db),
	}
}

func (h *NewsHandler) GetAllNews(c *gin.Context) {
	now := time.Now()
	news := h.newsRepo.GetNews(now, c)

	if len(news) == 0 {
		yesterday := now.AddDate(0, 0, -1)
		news = h.newsRepo.GetNews(yesterday, c)
	}

	slice.Shuffle(news)
	newsResponses := models.ToGetAllNewsResponse(news)

	c.JSON(http.StatusOK, gin.H{
		"news": newsResponses,
	})
}

func (h *NewsHandler) GetNewsReactionsById(c *gin.Context) {
	newsId := c.Param("news-id")
	reactions, err := h.reactionRepo.GetNewsReactionsById(newsId, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reactions": reactions,
	})
}

func (h *NewsHandler) Scrape(c *gin.Context) {
	scraper := scraping.NewSiteScraping()

	go func() {
		news := h.scrapeNews(scraper)

		h.scrapeReactions(news, scraper)

		f := &cache.CDNInvalidatorFactory{}
		inv := f.CreateInvalidator()
		inv.Invalidate()
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News scraping has been started in the background",
	})
}

func (h *NewsHandler) scrapeNews(scraper *scraping.SiteScraping) []models.News {
	news, err := scraper.ScrapeNews()
	if err != nil {
		log.Printf("Error scraping news: %v", err)
		return nil
	}

	return news
}

func (h *NewsHandler) scrapeReactions(news []models.News, scraper *scraping.SiteScraping) []models.News {
	for i := range news {
		reactions, err := scraper.ScrapeReactions(news[i])
		if err != nil {
			log.Printf("Error scraping reactions: %v", err)
		}
		news[i].Reactions = reactions
		h.db.Create(&news[i])
	}

	log.Printf("News scraping completed successfully, articles saved: %d", len(news))

	return news
}
