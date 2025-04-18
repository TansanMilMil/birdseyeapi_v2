package models

import (
	"time"

	"gorm.io/gorm"
)

// NewsReaction represents a reaction to a news article
type NewsReaction struct {
	gorm.Model
	ID           uint      `gorm:"primarykey"`
	NewsID       uint      `json:"newsId"`
	ReactionType string    `json:"reactionType"`
	Count        int       `json:"count"`
	Source       string    `json:"source"`
	ScrapedAt    time.Time `json:"scrapedAt"`
}