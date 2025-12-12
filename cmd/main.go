package main

import (
	"log"

	"erp-project/database"
	"erp-project/handlers"
	"erp-project/middleware"
	"erp-project/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.DB.Close()

	// Initialize repositories
	productRepo := repositories.NewProductRepository(database.DB)
	customerRepo := repositories.NewCustomerRepository(database.DB)
	orderRepo := repositories.NewOrderRepository(database.DB)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productRepo)
	customerHandler := handlers.NewCustomerHandler(customerRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo, productRepo, customerRepo)

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())

	// Product routes
	products := r.Group("/api/products")
	{
		products.POST("/", productHandler.CreateProduct)
		products.GET("/", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Customer routes
	customers := r.Group("/api/customers")
	{
		customers.POST("/", customerHandler.CreateCustomer)
		customers.GET("/", customerHandler.GetAllCustomers)
		customers.GET("/:id", customerHandler.GetCustomerByID)
	}

	// Order routes
	orders := r.Group("/api/orders")
	{
		orders.POST("/", orderHandler.CreateOrder)
		orders.GET("/", orderHandler.GetOrders)
		orders.GET("/:id/items", orderHandler.GetOrderItems)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "connected",
			"version":  "1.0.0",
		})
	})

	// Start server
	log.Println("Starting ERP server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
