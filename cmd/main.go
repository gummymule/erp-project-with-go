package main

import (
	"log"
	"os"

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

	// Add a root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ERP API is running",
			"version": "1.0.0",
			"routes": gin.H{
				"products":  "/api/products",
				"customers": "/api/customers",
				"orders":    "/api/orders",
				"health":    "/health",
			},
		})
	})

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

	// Health check (already correct)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "connected",
			"version":  "1.0.0",
		})
	})

	// Get port from Railway environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	// Start server
	log.Printf("Starting ERP server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
