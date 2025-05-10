package main

import (
	"log"

	api "github.com/birdseyeapi/birdseyeapi_v2/go/src/api"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/db"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	api.Init(db)
}
