package repository

import (
	"net/http"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

func (n *NewsRepository) GetNews(t time.Time, c *gin.Context) []models.News {
	var news []models.News

	result := n.db.
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
