package models

import (
	"time"

	"gorm.io/gorm"
)

// News represents a news article
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

// NewsWithReactionCount represents a news article with reaction count
type NewsWithReactionCount struct {
	News
	ReactionCount int `json:"reactionCount"`
}

// NewsSummarizeRequest represents a request to summarize a news article
type NewsSummarizeRequest struct {
	URL string `json:"url"`
}