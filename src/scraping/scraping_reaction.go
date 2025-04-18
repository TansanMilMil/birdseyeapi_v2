package scraping

import (
	"github.com/birdseyeapi/birdseyeapi_v2/internal/models"
)

// ScrapingReaction defines the interface for reaction scrapers
type ScrapingReaction interface {
	// ExtractReactions extracts reactions for a specific article
	ExtractReactions(articleURL, title string) ([]models.NewsReaction, error)
	
	// GetSourceBy returns the source identifier
	GetSourceBy() string
}