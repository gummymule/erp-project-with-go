package middleware

import (
	"erp-project/utils"
	"log"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get environment (development/production)
				env := os.Getenv("GO_ENV")

				// Log details
				log.Printf("⚠️ [PANIC RECOVERED]")
				log.Printf("   Error: %v", err)
				log.Printf("   Method: %s", c.Request.Method)
				log.Printf("   Path: %s", c.Request.URL.Path)
				log.Printf("   Client IP: %s", c.ClientIP())

				// Include stack trace in logs
				stack := debug.Stack()
				log.Printf("   Stack Trace:\n%s", string(stack))

				// Prepare response data
				var responseData interface{}
				if env == "development" || env == "" {
					// In development, include more details
					responseData = map[string]interface{}{
						"error":   err,
						"path":    c.Request.URL.Path,
						"method":  c.Request.Method,
						"message": "Internal server error (development mode)",
					}
				} else {
					// In production, generic message
					responseData = "An unexpected error occurred"
				}

				// Send response using new format
				utils.InternalErrorResponse(c,
					"Internal Server Error",
					responseData,
				)

				c.Abort()
			}
		}()

		c.Next()
	}
}
