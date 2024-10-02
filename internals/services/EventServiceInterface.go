package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/types"
	"github.com/lordofthemind/EventureGo/internals/utils"
)

// EventServiceInterface defines the methods required for managing events
type EventServiceInterface interface {
	// CreateEventService handles the creation of a new event
	CreateEventService(ctx context.Context, OrganizerID uuid.UUID, event *utils.EventDTO) (*types.EventType, error)

	// // FindEventByIDService retrieves an event by its unique identifier
	// FindEventByIDService(ctx context.Context, eventID uuid.UUID) (*types.EventType, error)

	// // FindEventsByOrganizerIDService retrieves all events created by a specific organizer
	// FindEventsByOrganizerIDService(ctx context.Context, organizerID uuid.UUID) ([]*types.EventType, error)

	// // UpdateEventService updates an existing event with new data
	// UpdateEventService(ctx context.Context, event *types.EventType) error

	// // DeleteEventByIDService removes an event from the system by its unique identifier
	// DeleteEventByIDService(ctx context.Context, eventID uuid.UUID) error

	// // FindAllEventsService retrieves all events in the system
	// FindAllEventsService(ctx context.Context) ([]*types.EventType, error)

	// // FindEventsByDateRangeService retrieves events occurring within a specified date range
	// FindEventsByDateRangeService(ctx context.Context, startDate, endDate time.Time) ([]*types.EventType, error)

	// // FindEventsByLocationService retrieves events occurring at a specific location
	// FindEventsByLocationService(ctx context.Context, location string) ([]*types.EventType, error)

	// // FindEventsByMultipleTagsService retrieves events matching specified tags
	// FindEventsByMultipleTagsService(ctx context.Context, tags []string) ([]*types.EventType, error)

	// // ActivateEventService activates an event, making it live
	// ActivateEventService(ctx context.Context, eventID uuid.UUID) error

	// // DeactivateEventService deactivates an event, making it inactive
	// DeactivateEventService(ctx context.Context, eventID uuid.UUID) error

	// // FindUpcomingEventsService retrieves events scheduled for future dates
	// FindUpcomingEventsService(ctx context.Context) ([]*types.EventType, error)

	// // FindPastEventsService retrieves events that have already occurred
	// FindPastEventsService(ctx context.Context) ([]*types.EventType, error)

	// // FindActiveEventsService retrieves all currently active events
	// FindActiveEventsService(ctx context.Context) ([]*types.EventType, error)

	// // FindInactiveEventsService retrieves all inactive events
	// FindInactiveEventsService(ctx context.Context) ([]*types.EventType, error)

	// // SearchEventsByTitleService searches for events by their title
	// SearchEventsByTitleService(ctx context.Context, title string) ([]*types.EventType, error)

	// // CancelEventService cancels an event, marking it as canceled
	// CancelEventService(ctx context.Context, eventID uuid.UUID) error

	// // RescheduleEventService reschedules an event to a new time slot
	// RescheduleEventService(ctx context.Context, eventID uuid.UUID, newStartTime, newEndTime time.Time) error

	// // CountTotalEventsService returns the total number of events in the system
	// CountTotalEventsService(ctx context.Context) (int64, error)

	// // CountEventsByOrganizerIDService returns the count of events organized by a specific user
	// CountEventsByOrganizerIDService(ctx context.Context, organizerID uuid.UUID) (int64, error)
}
