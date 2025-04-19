package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	newsHandler := NewNewsHandler(db)

	newsRoutes := r.Group("/news")
	{
		newsRoutes.GET("/", newsHandler.GetAllNews)
		newsRoutes.GET("/:id", newsHandler.GetNewsById)
		newsRoutes.POST("/", newsHandler.CreateNews)
		newsRoutes.POST("/scrape", newsHandler.ScrapeNews)
		newsRoutes.POST("/summarize", newsHandler.SummarizeNews)
	}
}
