package db

import "github.com/birdseyeapi/birdseyeapi_v2/go/src/models"

func GetMigrationModels() []interface{} {
	return []interface{}{
		&models.News{},
		&models.NewsReaction{},
	}
}
