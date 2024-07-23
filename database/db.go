package database

import (
	//	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"product-info/database/model"
	log "github.com/sirupsen/logrus"
)

/*
type Product struct {
	Id          int     `json:"id"		gorm:"primaryKey;type:int;autoIncrement"`
	Name        string  `json:"name"	gorm:"type:varcha(64)"`
	Description string  `json:"description"	gorm:"type:varchar(256)"`
	Price       float32 `json:"price"	gorm:"type:numeric"` 			//UnitPrice
	Quantity    int     `json:"quantity"	gorm:"type:int"`
	Discount    int     `json:"discount"	gorm:"type:int"` 			//MaxDiscountPercent
	Country     string  `json:"country"	gorm:"type:varchar(64)"`				
	Region      string  `json:"region"	gorm:"type:varchar(8)"`	
}
*/

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
