package reaction

import (
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
)

type ScrapingReaction interface {
	ExtractReactions(newsId uint, articleURL string, title string) ([]models.NewsReaction, error)
	GetSourceBy() string
}
