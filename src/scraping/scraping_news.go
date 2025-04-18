package scraping

import (
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

// ScrapingNews defines the interface for news scrapers
type ScrapingNews interface {
	// ExtractNews extracts news from a specific source
	ExtractNews() ([]models.News, error)
	
	// GetSourceBy returns the source identifier
	GetSourceBy() string
}