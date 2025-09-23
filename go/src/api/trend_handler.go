package api

import (
	"log"
	"net/http"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/trend"
	"github.com/gin-gonic/gin"
)

type TrendHandler struct {
	trendCatcher trend.TrendCatcher
}

func NewTrendHandler() *TrendHandler {
	fac := &trend.TrendCatcherFactory{}
	return &TrendHandler{
		trendCatcher: fac.CreateTrendCatcher(),
	}
}

func (h *TrendHandler) GetTrends(c *gin.Context) {
	news, err := h.trendCatcher.GetTrends()
	if err != nil {
		log.Printf("Error getting trends: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve trends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trends": news,
	})
}
