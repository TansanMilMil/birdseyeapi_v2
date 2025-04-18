package db

import (
	"birdseyeapi_v2/src/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	// Get database connection parameters from environment variables
	// or use default values if not provided
	username := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "root1603")
	host := getEnv("DB_HOST", "mysql")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "birdseyeapi")

	// Create the database connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		username, password, host, port, dbname,
	)

	// Connect to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Auto-migrate the database schema
	err = db.AutoMigrate(&models.News{}, &models.NewsReaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getEnv gets an environment variable value or returns a default value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}