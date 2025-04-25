package models

import (
	"time"

	"gorm.io/gorm"
)

type NewsReaction struct {
	gorm.Model
	ID              uint      `gorm:"primarykey"`
	NewsID          uint      `json:"newsId"`
	Author          string    `json:"author"`
	Comment         string    `json:"comment"`
	ScrapedDateTime time.Time `json:"scrapedDateTime"`
	CommentUrl      string    `json:"commentUrl"`
}
