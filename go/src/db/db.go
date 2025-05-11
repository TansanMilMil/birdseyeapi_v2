package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	password := os.Getenv("MYSQL_ROOT_PASSWORD")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		Username,
		password,
		Host,
		Port,
		DBName,
		Charset,
		ParseTime,
		Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(GetMigrationModels()...)
	if err != nil {
		return nil, err
	}

	return db, nil
}
