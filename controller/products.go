package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"product-info/rmqsender"
	"product-info/handler"
	"product-info/database/model"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var products []*model.Product

var (
	IdAccessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_id_access_total",
			Help: "Total number of times product ids are accessed",
		},
		[]string{"id"},
	)
	numberOfUpdateRequests = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "update_requests_total",
			Help: "Product Update requested.",
		},
	)
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests",
		},
		[]string{"path"},
	)
	TotalHTTPMethods = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_method_total",
			Help: "Number of request by http method",
		},
		[]string{"method"},
	)
	HttpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time_seconds", 
			Help: "Duration of HTTP requests",
		},
		[]string{"path"},
	)
)


// ListProducts - Get all Products Information
//
//	@Summary 	List all Products information from the server
//	@Description 	Get method for listing all the products from server
//	@Tags		Products
//	@Produce	json
//	@Accept		json
//	@Success	200
//	@Failure	404
//	@Router		/products	[get]
func ListProducts(c *gin.Context) {
	
	p := handler.Product{}
	p.Build(nil)

	response, status := p.GetProducts()

	c.IndentedJSON(status, response)
	return
}

// GetProduct - Get information about particular Product
//
//	@Summary 	Get particular Product's information from the server
//	@Description 	Get method for particular product information from server
//	@Tags		Product
//	@Param		id 	path 	int 	true 	"Product ID"
//	@Produce	json
//	@Accept		json
//	@Success	200
//	@Failure	404
//	@Router		/product/{id}	[get]
func GetProduct(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	
	p := handler.Product{}
	p.Build(nil)

	response, status := p.GetProduct(id)

	c.IndentedJSON(status, response)
	return
}


// DeleteProduct - Delete existing product from Product information system
//
//	@Summary 	Delete Product from product information system
//	@Description 	Delete method for particular product's information on the server
//	@Tags		Product
//	@Param		id 	path 	int 	true 	"Product ID"
//	@Produce	json
//	@Accept		json
//	@Success	200
//	@Failure	404
//	@Failure	400
//	@Router		/product/{id}	[delete]
func DeleteProduct(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	//idInt, _ := strconv.Atoi(id)

	p := handler.Product{}
	p.Build(nil)

	status := p.DeleteProduct(id)
	response := gin.H{"message": "successfully delete product"}

	if status != http.StatusOK {
		c.IndentedJSON(status, gin.H{"message":"failed to delete product"})
		return
	}

	c.IndentedJSON(status, response)
	return
}



// PostProduct - add new product info exist product information system
//
//	@Summary 	add a new Product into product information system
//	@Description 	create method for adding new product info product information system on the server
//	@Tags		Product
//	@Param		addProduct 	body 	model.Product 	true 	"Product to add Product information system"
//	@Produce	json
//	@Accept		json
//	@Success	200
//	@Failure	404
//	@Failure	400
//	@Router		/product	[post]
func PostProduct(c *gin.Context) {
	//var cproduct model.CreateProduct
	var cproduct model.Product

	if err := c.BindJSON(&cproduct); err != nil {
		log.Errorf("Error binding at Product : %+v ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error binding product"})
		return
	}

	p := handler.Product{}
	p.Build(nil)
	response, status := p.CreateProduct(&cproduct)
	c.IndentedJSON(status, response)
}


// UpdateProduct - Update existing product's information
//
//	@Summary 	Update Product's information
//	@Description 	Update method for particular product's information on the server
//	@Tags		Product
//	@Param		id 	path 	int 	true 	"Product ID"
//	@Param		updateProduct 	body 	model.Product 	true 	"Product details to update"
//	@Produce	json
//	@Accept		json
//	@Success	200
//	@Failure	404
//	@Failure	400
//	@Router		/product/{id}	[patch]
func UpdateProduct(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	var product model.Product

	if err := c.BindJSON(&product); err != nil {
		log.Errorf("Error binding at Product : %+v ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error binding product"})
		return
	}

	numberOfUpdateRequests.Inc()
	IdAccessCounter.WithLabelValues(id).Inc()
	idInt, _ := strconv.Atoi(id)

	if product.Id != idInt {
		log.Errorf("mis-match in product id btw url path %v and body %v", idInt, product.Id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "mis-match in product id btw url path and body"})
		return
	}

	p := handler.Product{}
	p.Build(nil)

	response, status := p.UpdateProduct(&product)

	fmt.Println("Status : ", status)
	fmt.Println("Respone : ", response)
			
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Json marshalling failed %+v ", err)
	}

	if status == http.StatusOK && response.Region == "EU" {
		senderEU := rmqsender.Connect("ProducedEU")
		defer senderEU.Close()
		senderEU.SendMessage(context.Background(), string(responseJson))
		c.IndentedJSON(http.StatusOK, response)
		return
	}
	
	if status == http.StatusOK && response.Region == "US" {
		senderUS := rmqsender.Connect("ProducedUS")
		defer senderUS.Close()
		senderUS.SendMessage(context.Background(), string(responseJson))
		c.IndentedJSON(http.StatusOK, response)
		return
	}

	c.JSON(status, response)
	return
}
