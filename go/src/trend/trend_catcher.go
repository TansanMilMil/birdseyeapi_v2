package trend

import "github.com/birdseyeapi/birdseyeapi_v2/go/src/models"

type TrendCatcher interface {
	GetTrends() ([]models.News, error)
}
