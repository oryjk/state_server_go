package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// InitDB initializes the database connection pool
func InitDB(dsn string) *gorm.DB {
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
