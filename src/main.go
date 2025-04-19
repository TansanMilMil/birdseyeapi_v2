package main

import (
	"log"

	"github.com/birdseyeapi/birdseyeapi_v2/src/api"
	db "github.com/birdseyeapi/birdseyeapi_v2/src/db"
	"github.com/gin-gonic/gin"
)

const (
	Port = "8080"
)

func main() {
	// Set up Gin
	r := gin.Default()
	
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	api.RegisterRoutes(r, db)

	log.Printf("Server starting on port %s", Port)
	if err := r.Run(":" + Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}