package responses

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// StandardResponse defines the structure for API responses
type StandardResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"requestId,omitempty"`
}

// NewGinResponse returns a gin standardized response
func NewGinResponse(c *gin.Context, status int, message string, data interface{}, err interface{}) StandardResponse {
	// Retrieve request ID from context, if available
	requestID, exists := c.Get("RequestID")
	if !exists {
		requestID = ""
	}

	return StandardResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Error:     err,
		Timestamp: time.Now().Format(time.RFC3339),
		RequestID: requestID.(string),
	}
}

// NewResponse returns a standardized response and includes request ID from context
func NewFiberResponse(c *fiber.Ctx, status int, message string, data interface{}, err interface{}) StandardResponse {
	// Get the request ID from the context with a default value if it doesn't exist
	requestID := c.Get("RequestID", "")

	return StandardResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Error:     err,
		Timestamp: time.Now().Format(time.RFC3339),
		RequestID: requestID,
	}
}
