package db

import (
	"fmt"
	"os"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	username := "root"
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	host := "mysql"
	port := "3306"
	dbname := "birds_eye"

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		username, password, host, port, dbname,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.News{}, &models.NewsReaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
