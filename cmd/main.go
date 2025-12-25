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

	// ✅ CORRECTED: NO trailing slash in groups!
	products := r.Group("/api/products") // NO trailing slash
	{
		products.POST("", productHandler.CreateProduct)      // POST /api/products
		products.GET("", productHandler.GetAllProducts)      // GET /api/products
		products.GET("list", productHandler.GetListProducts) // GET /api/products/list
		products.GET(":id", productHandler.GetProductByID)   // GET /api/products/:id
		products.PUT(":id", productHandler.UpdateProduct)    // PUT /api/products/:id
		products.DELETE(":id", productHandler.DeleteProduct) // DELETE /api/products/:id
	}

	// ✅ CORRECTED: NO trailing slash
	customers := r.Group("/api/customers")
	{
		customers.POST("", customerHandler.CreateCustomer)    // POST /api/customers
		customers.GET("", customerHandler.GetAllCustomers)    // GET /api/customers
		customers.GET(":id", customerHandler.GetCustomerByID) // GET /api/customers/:id
	}

	// ✅ CORRECTED: NO trailing slash
	orders := r.Group("/api/orders")
	{
		orders.POST("", orderHandler.CreateOrder)           // POST /api/orders
		orders.GET("", orderHandler.GetOrders)              // GET /api/orders
		orders.GET(":id/items", orderHandler.GetOrderItems) // GET /api/orders/:id/items
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
