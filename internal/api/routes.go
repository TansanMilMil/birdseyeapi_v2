package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers all API endpoints
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Create a new instance of NewsHandler with the provided db
	newsHandler := NewNewsHandler(db)

	// Create an API group
	api := r.Group("/api")
	{
		// News endpoints
		newsRoutes := api.Group("/news")
		{
			newsRoutes.GET("/", newsHandler.GetAllNews)
			newsRoutes.GET("/:id", newsHandler.GetNewsById)
			newsRoutes.POST("/", newsHandler.CreateNews)
			newsRoutes.POST("/scrape", newsHandler.ScrapeNews)
			newsRoutes.POST("/summarize", newsHandler.SummarizeNews)
		}
	}
}