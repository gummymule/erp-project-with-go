package middleware

import (
	"bytes"
	"io"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// start timer
		start := time.Now()

		// read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// create custom response writer
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// process request
		c.Next()

		// log details
		duration := time.Since(start)

		log.Printf("[%s] %s %s - %d - %v",
			c.Request.Method,
			c.Request.URL.Path,
			string(requestBody),
			c.Writer.Status(),
			duration,
		)

		// log response for errors
		if c.Writer.Status() >= 400 {
			log.Printf("Error Response: %s", blw.body.String())
		}

	}
}
