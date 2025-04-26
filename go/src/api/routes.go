package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	newsHandler := NewNewsHandler(db)

	newsRoutes := r.Group("/news")
	{
		newsRoutes.GET("/today-news", newsHandler.GetAllNews)
		newsRoutes.GET("/news-reactions/:news-id", newsHandler.GetNewsReactionsById)
		newsRoutes.POST("/scrape", newsHandler.ScrapeNews)
	}
}
