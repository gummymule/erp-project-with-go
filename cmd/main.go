package main

import (
	"log"
	"os"
	"time"

	"erp-project/database"
	"erp-project/handlers"
	"erp-project/middleware"
	"erp-project/repositories"
	"erp-project/utils"

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

	// ✅ Add recovery middleware FIRST (to catch any panics)
	r.Use(middleware.Recovery())

	// ✅ Add your existing middleware - REMOVE the duplicate logging middleware below
	r.Use(middleware.RequestLogger())

	// Add a root route with standardized response format
	r.GET("/", func(c *gin.Context) {
		utils.SuccessResponse(c, "ERP API is running", map[string]interface{}{
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
			"endpoints": map[string]interface{}{
				"products": map[string]string{
					"create":   "POST /api/products",
					"get_all":  "GET /api/products",
					"get_list": "GET /api/products/list",
					"get_one":  "GET /api/products/:id",
					"update":   "PUT /api/products/:id",
					"delete":   "DELETE /api/products/:id",
				},
				"customers": map[string]string{
					"create":  "POST /api/customers",
					"get_all": "GET /api/customers",
					"get_one": "GET /api/customers/:id",
				},
				"orders": map[string]string{
					"create":    "POST /api/orders",
					"get_all":   "GET /api/orders",
					"get_items": "GET /api/orders/:id/items",
				},
				"health":   "GET /health",
				"debug_db": "GET /debug/db",
			},
		})
	})

	// ✅ FIXED: Add leading slashes to all routes in groups
	products := r.Group("/api/products")
	{
		products.POST("/", productHandler.CreateProduct)      // ✅ POST /api/products/
		products.GET("/", productHandler.GetAllProducts)      // ✅ GET /api/products/
		products.GET("/list", productHandler.GetListProducts) // ✅ GET /api/products/list
		products.GET("/:id", productHandler.GetProductByID)   // ✅ GET /api/products/:id
		products.PUT("/:id", productHandler.UpdateProduct)    // ✅ PUT /api/products/:id
		products.DELETE("/:id", productHandler.DeleteProduct) // ✅ DELETE /api/products/:id
	}

	// Customer routes - FIXED with leading slashes
	customers := r.Group("/api/customers")
	{
		customers.POST("/", customerHandler.CreateCustomer)    // ✅ POST /api/customers/
		customers.GET("/", customerHandler.GetAllCustomers)    // ✅ GET /api/customers/
		customers.GET("/:id", customerHandler.GetCustomerByID) // ✅ GET /api/customers/:id
	}

	// Order routes - FIXED with leading slashes
	orders := r.Group("/api/orders")
	{
		orders.POST("/", orderHandler.CreateOrder)           // ✅ POST /api/orders/
		orders.GET("/", orderHandler.GetOrders)              // ✅ GET /api/orders/
		orders.GET("/:id/items", orderHandler.GetOrderItems) // ✅ GET /api/orders/:id/items
	}

	// Health check with standardized response format
	r.GET("/health", func(c *gin.Context) {
		// Check database connection
		err := database.DB.Ping()
		databaseStatus := "connected"
		if err != nil {
			databaseStatus = "disconnected"
		}

		utils.SuccessResponse(c, "System is healthy", map[string]interface{}{
			"status":    "OK",
			"database":  databaseStatus,
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Debug endpoint to test database with standardized response format
	r.GET("/debug/db", func(c *gin.Context) {
		var productCount, customerCount, orderCount int
		var errorMsg string

		// Get counts with error handling
		if err := database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount); err != nil {
			errorMsg = "Failed to count products: " + err.Error()
			productCount = -1
		}
		if err := database.DB.QueryRow("SELECT COUNT(*) FROM customers").Scan(&customerCount); err != nil {
			errorMsg = "Failed to count customers: " + err.Error()
			customerCount = -1
		}
		if err := database.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orderCount); err != nil {
			errorMsg = "Failed to count orders: " + err.Error()
			orderCount = -1
		}

		// Check if there were any errors
		if errorMsg != "" {
			utils.InternalErrorResponse(c, "Database statistics error", errorMsg)
			return
		}

		utils.SuccessResponse(c, "Database statistics", map[string]interface{}{
			"database":        "connected",
			"products_count":  productCount,
			"customers_count": customerCount,
			"orders_count":    orderCount,
			"total_records":   productCount + customerCount + orderCount,
			"timestamp":       time.Now().Format(time.RFC3339),
		})
	})

	// Test endpoint with standardized response format
	r.POST("/api/test-simple", func(c *gin.Context) {
		log.Println("🎯 SIMPLE POST ENDPOINT HIT!")

		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ValidationErrorResponse(c, "Invalid request data", err.Error())
			return
		}

		log.Printf("📦 Received data: %v", data)
		utils.CreatedResponse(c, "Simple POST works!", map[string]interface{}{
			"received_data": data,
			"timestamp":     time.Now().Format(time.RFC3339),
		})
	})

	// Get port from Railway environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	// Start server
	log.Printf("🚀 Starting ERP server on :%s", port)
	log.Printf("📊 Available at: http://localhost:%s", port)
	log.Printf("📋 Root endpoint: http://localhost:%s/", port)
	log.Printf("🏥 Health check: http://localhost:%s/health", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
