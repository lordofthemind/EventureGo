package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/types"
)

// RegisterEventRequest defines the structure for registering a new event
type RegisterEventRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=100"`        // Required event title
	Description string    `json:"description" validate:"max=1000"`                // Optional description with a limit
	StartTime   time.Time `json:"start_time" validate:"required"`                 // Required start time
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"` // Required end time must be greater than start time
	Location    string    `json:"location" validate:"required"`                   // Required location
	Tags        []string  `json:"tags" validate:"dive,required"`                  // Optional tags, can be empty
}

// EventDTO is the internal representation of the event data
type EventDTO struct {
	Title       string
	Description string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	Tags        []string
}

// TransformToEventDTO converts the incoming request to an EventDTO for internal use
func TransformToEventDTO(eventReq RegisterEventRequest) *EventDTO {
	return &EventDTO{
		Title:       eventReq.Title,
		Description: eventReq.Description,
		Location:    eventReq.Location,
		StartTime:   eventReq.StartTime,
		EndTime:     eventReq.EndTime,
		Tags:        eventReq.Tags,
	}
}

// RegisterEventResponse defines the structure for the response after registering a new event
type RegisterEventResponse struct {
	ID          uuid.UUID         `json:"id"`           // The unique identifier for the event
	Title       string            `json:"title"`        // Event title
	Description string            `json:"description"`  // Event description
	StartTime   time.Time         `json:"start_time"`   // Start time of the event
	EndTime     time.Time         `json:"end_time"`     // End time of the event
	Location    string            `json:"location"`     // Location of the event
	OrganizerID uuid.UUID         `json:"organizer_id"` // ID of the organizer
	Guests      []types.GuestType `json:"guests"`       // List of guests (can be empty)
	CreatedAt   time.Time         `json:"created_at"`   // Timestamp when the event was created
	UpdatedAt   time.Time         `json:"updated_at"`   // Timestamp when the event was last updated
	IsActive    bool              `json:"is_active"`    // Status of the event (active/inactive)
	Tags        []string          `json:"tags"`         // Tags associated with the event
}

// TransformToRegisterEventResponse converts the EventType to RegisterEventResponse
func TransformToRegisterEventResponse(event *types.EventType) *RegisterEventResponse {
	return &RegisterEventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Location:    event.Location,
		OrganizerID: event.OrganizerID,
		Guests:      event.Guests, // Ensure this is handled properly (not nil)
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		IsActive:    event.IsActive,
		Tags:        event.Tags,
	}
}
