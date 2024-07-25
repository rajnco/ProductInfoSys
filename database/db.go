package database

import (
	//	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"product-info/database/model"
	log "github.com/sirupsen/logrus"
)

var dbClient *gorm.DB
type product model.Product


func InitDB(dbname string) *gorm.DB {

	//db, err := gorm.Open(sqlite.Open("product.db"), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
	}

	dbClient = db

	return dbClient
}


func GetDB(dbname string) *gorm.DB {
	dbClient = InitDB(dbname)
	return dbClient
}


func StartMigration(dbname string) error {
	db := GetDB(dbname)
	if err := db.Table("products").AutoMigrate(&product{}); err != nil {
		log.Println("failed to create DB")
	}
	return nil
}


func DropAllTables(dbname string) error {
	db := GetDB(dbname)

	if err := db.Migrator().DropTable(); err != nil {
		return err 
	}
	log.Warnln("Delete old tables...")
	return nil
}
