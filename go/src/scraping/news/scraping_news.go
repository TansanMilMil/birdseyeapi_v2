package news

import (
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
)

type ScrapingNews interface {
	ExtractNews() ([]models.News, error)
	GetSourceBy() string
}
