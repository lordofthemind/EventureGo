package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
	"github.com/lordofthemind/EventureGo/internals/validators"
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
	// Retrieve OrganizerID from the context set by middleware
	organizerID, exists := c.Get("userID")
	if !exists {
		response := responses.NewGinResponse(c, http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Type assert the organizerID to uuid.UUID
	userID, ok := organizerID.(uuid.UUID)
	if !ok {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Invalid User ID type", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Log the user ID for audit purposes
	log.Println("User ID:", userID)

	// Parse the request body into RegisterEventRequest struct
	var eventRequest utils.RegisterEventRequest
	if err := c.ShouldBindJSON(&eventRequest); err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate the event request using a validator
	if validationErr := validators.ValidateEventRequest(eventRequest); validationErr != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Validation error", nil, validationErr.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Transform request data to internal DTO
	eventDTO := utils.TransformToEventDTO(eventRequest)

	// Call the service to create the event
	createdEvent, err := h.service.CreateEventService(c.Request.Context(), userID, eventDTO)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to create event", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Transform createdEvent to RegisterEventResponse and return it
	responseData := utils.TransformToRegisterEventResponse(createdEvent)

	// Return the standardized response
	response := responses.NewGinResponse(c, http.StatusOK, "Event created successfully", responseData, nil)
	c.JSON(http.StatusOK, response)
}
