package handler 

import (
	"product-info/database/model"
	"testing"
	"strconv"
	//"fmt"

	"github.com/stretchr/testify/assert"
)

func TestProducts(t *testing.T) {
	var product model.Product
	var products []model.Product

	// Create Product
	p := Product{}
	p.Build(nil)

	product = model.Product{Id: 0, Name: "Product 2", Description: "Table", Price: 100, Quantity: 100, Discount: 10, Country: "EU", Region: "EU"}

	response, status := p.CreateProduct(&product)

	assert.Equal(t, status, 201)
	assert.Equal(t, response.Name, "Product 2")

	// Get Products
	products, status = p.GetProducts()

	assert.Equal(t, status, 200)
	assert.Equal(t, products[0].Name, "Product 2")

	// Update Product
	product = model.Product{Id: response.Id, Name: "Product 3", Description: "Table", Price: 100, Quantity: 100, Discount: 10, Country: "EU", Region: "EU"}

	response, status = p.UpdateProduct(&product)

	assert.Equal(t, status, 200)
	assert.Equal(t, response.Name, "Product 3")

	// Get Product
	idStr := strconv.Itoa(response.Id) 

	response, status = p.GetProduct(idStr)

	assert.Equal(t, status, 200)
	assert.Equal(t, response.Name, "Product 3")


	// Delete Product
	status = p.DeleteProduct(idStr)

	assert.Equal(t, status, 200)


	// Clean DB
	err := p.Clean()

	assert.Nil(t, err)

	// Total Count
	count := p.TotalCount()

	assert.Equal(t, int64(0), count)

}
