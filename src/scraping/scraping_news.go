package scraping

import (
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

type ScrapingNews interface {
	ExtractNews() ([]models.News, error)
	GetSourceBy() string
}