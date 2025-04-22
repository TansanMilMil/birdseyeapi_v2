package reaction

import (
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

type ScrapingReaction interface {
	ExtractReactions(articleURL, title string) ([]models.NewsReaction, error)
	GetSourceBy() string
}
