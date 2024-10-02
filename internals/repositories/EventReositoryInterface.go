package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/types"
)

// EventRepositoryInterface defines the methods for handling events in the system
type EventRepositoryInterface interface {
	// CreateEvent creates a new event
	CreateEvent(ctx context.Context, event *types.EventType) (*types.EventType, error)

	// FindEventByID finds an event by its ID
	FindEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error)

	// FindEventsByOrganizerID retrieves events organized by a specific user
	FindEventsByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]*types.EventType, error)

	// UpdateEvent updates an existing event
	UpdateEvent(ctx context.Context, event *types.EventType) error

	// DeleteEventByID deletes an event by its ID
	DeleteEventByID(ctx context.Context, eventID uuid.UUID) error

	// FindAllEvents retrieves all events in the system
	FindAllEvents(ctx context.Context) ([]*types.EventType, error)

	// FindEventsByDateRange finds events occurring between a start and end date
	FindEventsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*types.EventType, error)

	// FindEventsByLocation finds events based on location
	FindEventsByLocation(ctx context.Context, location string) ([]*types.EventType, error)

	// FindEventsByMultipleTags retrieves events that match any of the specified tags
	FindEventsByMultipleTags(ctx context.Context, tags []string) ([]*types.EventType, error)

	// ActivateEvent activates an event, making it live
	ActivateEvent(ctx context.Context, eventID uuid.UUID) error

	// DeactivateEvent deactivates an event, making it inactive
	DeactivateEvent(ctx context.Context, eventID uuid.UUID) error

	// FindUpcomingEvents retrieves events scheduled for future dates
	FindUpcomingEvents(ctx context.Context) ([]*types.EventType, error)

	// FindPastEvents retrieves events that have already occurred
	FindPastEvents(ctx context.Context) ([]*types.EventType, error)

	// FindActiveEvents retrieves all currently active events
	FindActiveEvents(ctx context.Context) ([]*types.EventType, error)

	// FindInactiveEvents retrieves all inactive events
	FindInactiveEvents(ctx context.Context) ([]*types.EventType, error)

	// SearchEventsByTitle finds events by searching their titles
	SearchEventsByTitle(ctx context.Context, title string) ([]*types.EventType, error)

	// CancelEvent cancels an event
	CancelEvent(ctx context.Context, eventID uuid.UUID) error

	// RescheduleEvent reschedules an event to a new date and time
	RescheduleEvent(ctx context.Context, eventID uuid.UUID, newStartTime, newEndTime time.Time) error

	// CountTotalEvents returns the total number of events in the system
	CountTotalEvents(ctx context.Context) (int64, error)

	// CountEventsByOrganizerID counts the number of events organized by a specific user
	CountEventsByOrganizerID(ctx context.Context, organizerID uuid.UUID) (int64, error)
}
