package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
)

type EventGinHandler struct {
	service services.EventServiceInterface
}

func NewEventGinHandler(service services.EventServiceInterface) *EventGinHandler {
	return &EventGinHandler{
		service: service,
	}
}

func (h *EventGinHandler) CreateEventHandler(c *gin.Context) {
	// Retrieve OrganizerID from the context set by the middleware
	organizerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assert the organizerID to uuid.UUID
	userID, ok := organizerID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User ID type"})
		return
	}

	// Log the user ID
	log.Println("User ID:", userID)

	// Parse the request body into RegisterEventRequest struct
	var eventRequest utils.RegisterEventRequest
	if err := c.ShouldBindJSON(&eventRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the event
	createdEvent, err := h.service.CreateEventService(c.Request.Context(), userID, &eventRequest) // Pass userID directly
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created event in the response
	c.JSON(http.StatusOK, createdEvent)
}
