package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var port = os.Getenv("GO_API_PORT")

func Init(db *gorm.DB) {
	r := gin.Default()

	r.Use(setConfig)

	RegisterRoutes(r, db)

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setConfig(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
