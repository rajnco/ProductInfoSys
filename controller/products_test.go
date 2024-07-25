package controller 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	"product-info/handler"
	"product-info/database/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)



func startTestServer() *gin.Engine {
	r := gin.Default()

	mainGroup := r.Group("/")
        {
                productsGroup := mainGroup.Group("/products")
                {
                        productsGroup.GET("", ListProducts)
                }
                
		productGroup := mainGroup.Group("/product")
                {
                        productGroup.POST("", PostProduct)
                        productGroup.GET("/:id", GetProduct)
                        productGroup.PATCH("/:id", UpdateProduct)
                        productGroup.DELETE("/:id", DeleteProduct)
                }
        }
	return r
}

func cleanProducts(t *testing.T){
	p := handler.Product{}
	p.Build(nil)

	if err := p.Clean(); err != nil {
		t.Error("Error clean up products table")
		t.Error(err)
	}
}


func TestProucts(t *testing.T){
	cleanProducts(t)

	r := startTestServer()
	w := httptest.NewRecorder()

	jsonProduct := []byte(`{"id": 0, "name": "Product 2", "description": "Table - hand made wooden table - US", 
				"price": 100, "quantity": 100, "discount": 10, "country": "EU", "region": "EU"}`)

	
	// POST
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonProduct))
	r.ServeHTTP(w, req)

	response := model.Product{}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Error : %v", err)
	}

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "Product 2", response.Name)


	// GET
	w2 := httptest.NewRecorder()
	url2 := fmt.Sprintf("/product/%s", strconv.Itoa(response.Id))
	req2, _ := http.NewRequest("GET", url2, nil)
	r.ServeHTTP(w2, req2)
	response2 := model.Product{}

	if err := json.Unmarshal(w2.Body.Bytes(), &response2); err != nil {
		t.Errorf("Error : %v", err)
	}

	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, "Product 2", response2.Name)

	// GET ALL
	w4 := httptest.NewRecorder()
	url4 := fmt.Sprintf("/products")
	req4, _ := http.NewRequest("GET", url4, nil)
	r.ServeHTTP(w4, req4)
	response4 := []model.Product{}

	if err := json.Unmarshal(w4.Body.Bytes(), &response4); err != nil {
		t.Errorf("Error : %v", err)
	}

	assert.Equal(t, 200, w4.Code)
	assert.Equal(t, "Product 2", response4[0].Name)

	// PATCH
	product := model.Product{}
	product.Id = response.Id
	product.Name = "Product New"
	product.Description = "changed"
	product.Price = 100
	product.Quantity = 1000
	product.Discount = 10
	product.Country = "US"
	product.Region = "US"

	jsonProduct2, _ := json.Marshal(&product)
	payload := []byte(jsonProduct2)

	w5 := httptest.NewRecorder()
	url5 := fmt.Sprintf("/product/%s", strconv.Itoa(response.Id))
	req5, _ := http.NewRequest("PATCH", url5, bytes.NewBuffer(payload))
	r.ServeHTTP(w5, req5)
	response5 := model.Product{}

	if err := json.Unmarshal(w5.Body.Bytes(), &response5); err != nil {
		t.Errorf("Error : %v", err)
	}

	assert.Equal(t, 200, w5.Code)
	assert.Equal(t, "Product New", response5.Name)

	// DELETE
	w3 := httptest.NewRecorder()
	url3 := fmt.Sprintf("/product/%s", strconv.Itoa(response.Id))
	req3, _ := http.NewRequest("DELETE", url3, nil)
	r.ServeHTTP(w3, req3)

	assert.Equal(t, 200, w3.Code)
}
