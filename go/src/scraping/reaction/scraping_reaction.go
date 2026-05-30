package reaction

import (
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/tebeka/selenium"
)

type ScrapingReaction interface {
	ExtractReactions(driver selenium.WebDriver, newsId uint, articleURL string, title string) ([]models.NewsReaction, error)
	GetSourceBy() string
}
