package reaction

import (
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
)

type ScrapingReaction interface {
	ExtractReactions(articleURL, title string) ([]models.NewsReaction, error)
	GetSourceBy() string
}
