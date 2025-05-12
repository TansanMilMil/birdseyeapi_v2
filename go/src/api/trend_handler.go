package api

import (
	"net/http"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/trend"
	"github.com/gin-gonic/gin"
)

type TrendHandler struct{}

func (h *TrendHandler) GetTrends(c *gin.Context) {
	fac := &trend.TrendCatcherFactory{}
	tCatcher := fac.CreateTrendCatcher()
	news := tCatcher.GetTrends()

	c.JSON(http.StatusOK, news)
}
