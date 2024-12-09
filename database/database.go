package database

import (
	_ "database/sql"
	"github.com/abolfazlalz/scalesops/golang/constants"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func Database() *gorm.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = gorm.Open(postgres.Open(constants.Database().Url()))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.AutoMigrate(&ImageRequest{}, &Image{}); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}
	return db
}
