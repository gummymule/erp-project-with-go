package middleware

import (
	"erp-project/utils"

	"log"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				// log stack trace
				// debug.PrintStack()

				utils.ErrorResponse(c, utils.NewAppError(
					500,
					"Internal Server Error",
					"Something went wrong",
				))

				c.Abort()
			}
		}()

		c.Next()
	}
}
