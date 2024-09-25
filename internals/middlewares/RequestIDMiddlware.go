package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates a request ID and adds it to the context and response
func RequestIDGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if client sent a request ID
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate a new UUID if not provided by the client
			requestID = uuid.New().String()
		}

		// Add the request ID to the context
		c.Set("RequestID", requestID)

		// Add the request ID to the response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Continue processing
		c.Next()
	}
}
