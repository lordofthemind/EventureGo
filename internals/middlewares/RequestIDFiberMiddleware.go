package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestIDFiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if client sent a request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			// Generate a new UUID if not provided by the client
			requestID = uuid.New().String()

		}

		// Add the request ID to the context
		c.Locals("RequestID", requestID)

		// Add the request ID to the response header
		c.Set("X-Request-ID", requestID)

		// Continue to the next middleware or handler
		return c.Next()
	}
}
