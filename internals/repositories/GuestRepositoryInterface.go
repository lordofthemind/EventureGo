package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/types"
)

// GuestRepositoryInterface defines the methods for handling guests in the system
type GuestRepositoryInterface interface {
	// AddGuest adds a new guest to an event
	AddGuest(ctx context.Context, guest *types.GuestType) (*types.GuestType, error)

	// AddGuest adds a new guest to an event
	AddBulkGuest(ctx context.Context, guest []*types.GuestType) ([]*types.GuestType, error)

	// FindGuestsByEventID retrieves all guests for a given event
	FindGuestsByEventID(ctx context.Context, eventID uuid.UUID) ([]*types.GuestType, error)

	// FindGuestByID retrieves a guest by their ID
	FindGuestByID(ctx context.Context, guestID uuid.UUID) (*types.GuestType, error)

	// UpdateGuest updates the information of an existing guest
	UpdateGuest(ctx context.Context, guest *types.GuestType) error

	// DeleteGuestByID removes a guest by their ID
	DeleteGuestByID(ctx context.Context, guestID uuid.UUID) error

	// CountGuestsByEventID counts the number of guests for a given event
	CountGuestsByEventID(ctx context.Context, eventID uuid.UUID) (int64, error)
}
