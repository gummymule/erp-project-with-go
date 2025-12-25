package main

import (
	"log"
	"os"

	"erp-project/database"
	"erp-project/handlers"
	"erp-project/middleware"
	"erp-project/repositories"

	"github.com/gin-contrib/cors"
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

	// ✅ ADD CORS MIDDLEWARE
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Allow all origins for now
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// Add your existing middleware
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

	// ✅ FIXED: Product routes - Add trailing slash to group to avoid redirects
	products := r.Group("/api/products/") // Add trailing slash here
	{
		products.POST("", productHandler.CreateProduct)      // Empty string
		products.GET("", productHandler.GetAllProducts)      // Empty string
		products.GET("list", productHandler.GetListProducts) // Relative path
		products.GET(":id", productHandler.GetProductByID)   // Relative path
		products.PUT(":id", productHandler.UpdateProduct)    // Relative path
		products.DELETE(":id", productHandler.DeleteProduct) // Relative path
	}

	// ✅ FIXED: Customer routes
	customers := r.Group("/api/customers/")
	{
		customers.POST("", customerHandler.CreateCustomer)
		customers.GET("", customerHandler.GetAllCustomers)
		customers.GET(":id", customerHandler.GetCustomerByID)
	}

	// ✅ FIXED: Order routes
	orders := r.Group("/api/orders/")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetOrders)
		orders.GET(":id/items", orderHandler.GetOrderItems)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"database": "connected",
			"version":  "1.0.0",
		})
	})

	// Debug endpoint to test database
	r.GET("/debug/db", func(c *gin.Context) {
		var productCount, customerCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount)
		database.DB.QueryRow("SELECT COUNT(*) FROM customers").Scan(&customerCount)

		c.JSON(200, gin.H{
			"database":        "connected",
			"products_count":  productCount,
			"customers_count": customerCount,
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
