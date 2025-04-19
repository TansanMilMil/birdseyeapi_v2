package models

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	ID              uint      `gorm:"primarykey"`
	ScrapingUnitID  uint      `json:"scrapingUnitId"`
	Title           string    `gorm:"size:200" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	SummarizedText  string    `gorm:"type:text" json:"summarizedText"`
	SourceBy        string    `json:"sourceBy"`
	ScrapedUrl      string    `gorm:"type:text" json:"scrapedUrl"`
	ScrapedDateTime time.Time `json:"scrapedDateTime"`
	ArticleUrl      string    `gorm:"type:text" json:"articleUrl"`
	ArticleImageUrl string    `gorm:"type:text" json:"articleImageUrl"`
	Reactions       []NewsReaction `gorm:"foreignKey:NewsID" json:"reactions"`
}

type NewsWithReactionCount struct {
	News
	ReactionCount int `json:"reactionCount"`
}

type NewsSummarizeRequest struct {
	URL string `json:"url"`
}