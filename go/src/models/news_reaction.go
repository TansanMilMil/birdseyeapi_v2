package models

import (
	"time"

	"gorm.io/gorm"
)

type NewsReaction struct {
	gorm.Model
	ID           uint      `gorm:"primarykey"`
	NewsID       uint      `json:"newsId"`
	ReactionType string    `json:"reactionType"`
	Count        int       `json:"count"`
	Source       string    `json:"source"`
	ScrapedAt    time.Time `json:"scrapedAt"`
}