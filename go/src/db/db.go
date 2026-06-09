package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 起動レース対策のリトライ設定。
// MySQL の初回ボリューム初期化は数秒かかるため、接続成功まで指数バックオフで待つ。
const (
	maxConnectRetries = 30
	initialBackoff    = 1 * time.Second
	maxBackoff        = 10 * time.Second
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

	db, err := connectWithRetry(dsn)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(GetMigrationModels()...)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectWithRetry は MySQL が接続可能になるまで指数バックオフでリトライする。
// gorm.Open は実接続を遅延させるため、sql.DB.Ping で実際の疎通を確認する。
func connectWithRetry(dsn string) (*gorm.DB, error) {
	backoff := initialBackoff
	var lastErr error

	for attempt := 1; attempt <= maxConnectRetries; attempt++ {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			var sqlDB *sql.DB
			if sqlDB, err = db.DB(); err == nil {
				if err = sqlDB.Ping(); err == nil {
					return db, nil
				}
			}
		}

		lastErr = err
		log.Printf("DB接続待ち (attempt %d/%d): %v", attempt, maxConnectRetries, err)
		if attempt == maxConnectRetries {
			break
		}

		time.Sleep(backoff)
		if backoff *= 2; backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	return nil, fmt.Errorf("database に接続できませんでした (%d回リトライ): %w", maxConnectRetries, lastErr)
}
