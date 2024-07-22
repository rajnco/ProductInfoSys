package handler

import (
	//"fmt"
	"net/http"
	"strconv"

	"product-info/database"
	"product-info/database/model"

	"gorm.io/gorm"
	log "github.com/sirupsen/logrus"
)

// ProductsTableName - table name for Product model in database
const ProductsTableName = "products"

// Product - data struct to access Products
type Product struct {
	db *gorm.DB
}

// Build - connect/re-use db connection
func (p *Product) Build(db *gorm.DB) {
	if db == nil {
		p.db = database.GetDB()
	} else {
		p.db = db
	}
}

// TotalCount - returns total count from table
func (p *Product) TotalCount() int64 {
	var count int64 
	if err := p.db.Table(ProductsTableName).Count(&count).Error; err != nil {
		return 0
	}
	return count
}


// CreateProduct - create new product in database table products
//func(p *Product) CreateProduct(cproduct *model.CreateProduct) (response model.Product, status int) {	
func(p *Product) CreateProduct(cproduct *model.Product) (response model.Product, status int) {	
	tx := p.db.Begin()
	
	if err := p.db.Model(&model.Product{}).Create(&cproduct).Error; err != nil {
		log.Printf(" error creating new product : %v ", err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return
	}
	//p.db.Model(model.Product{Id: cproduct.Id}).First(&response)
	//if err := p.db.Model(&model.Product{Id: cproduct.Id}).First(&response).Error; err != nil {
	if err := p.db.Model(&model.Product{}).Last(&response).Error; err != nil {
		log.Printf(" error creating new product : %v ", err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return
	}
	tx.Commit()

	status = http.StatusCreated
	return
}

// UpdateProduct - updates product info in database
func (p *Product) UpdateProduct(product *model.Product) (response model.Product, status int) {
	tx := p.db.Begin()

	if err := p.db.Model(&model.Product{}).Where("id = ? ", product.Id).Updates(&product).Error; err != nil {
		log.Printf(" error creating new product : %v ", err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return
	}

	tx.Commit()
	
	// status = http.StatusOK
	// return
	return p.GetProduct(strconv.Itoa(product.Id))
}

// GetProduct - get particular product from products list
func (p *Product) GetProduct(id string) (response model.Product, status int) {
	//var response model.Product

	//if err := p.db.Model(&model.Product{}).Where("id = ?", id).Find(&response).Error; err != nil {
	if err := p.db.Where("id = ?", id).Find(&response).Error; err != nil {
		log.Printf("error retrive product from DB : %v", err)
		status = http.StatusInternalServerError
		return
	}
	status = http.StatusOK
	return
}

// DeleteProduct - delete a product ffrom products list
func (p *Product) DeleteProduct(id string) (status int) {

	var product model.Product

	tx := p.db.Begin()

	if err := p.db.Model(&model.Product{}).Where("id = ?", id).Delete(&product).Error; err != nil {
	//if err := p.db.Model(&model.Product{}).Clauses(clause.Returning{}).Where("id = ?", id).Delete(&product).Error; err != nil {
		log.Printf("error delete product from DB : %v", err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return
	}
	status = http.StatusOK
	//fmt.Println("Deleted : %v ", product)
	tx.Commit()
	return
}


// GetProducts - list all the products 
func (p *Product) GetProducts() (response []model.Product, status int) {
	if err := p.db.Find(&response).Error; err != nil {
		log.Printf("error getting all products from DB : %v", err)
		status = http.StatusInternalServerError
		return
	}
	status = http.StatusOK
	return
}


// Clean - clean the table - delete all the records
func (p *Product) Clean() error {
	result := p.db.Exec("DELETE FROM " + ProductsTableName)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
