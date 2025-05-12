package repository

import (
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsReactionRepository struct {
	db *gorm.DB
}

func NewNewsReactionRepository(db *gorm.DB) *NewsReactionRepository {
	return &NewsReactionRepository{
		db: db,
	}
}

func (repo *NewsReactionRepository) GetNewsReactionsById(newsId string, c *gin.Context) ([]models.NewsReaction, error) {
	var reactions []models.NewsReaction

	result := repo.db.
		Where("news_id = ?", newsId).
		Limit(100).
		Find(&reactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return reactions, nil
}
