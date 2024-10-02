package postgresdb

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"gorm.io/gorm"
)

type postgresGuestRepository struct {
	db *gorm.DB
}

// NewPostgresGuestRepository initializes a new instance of the guest repository.
func NewPostgresGuestRepository(db *gorm.DB) repositories.GuestRepositoryInterface {
	return &postgresGuestRepository{
		db: db,
	}
}

// AddGuest adds a new guest to an event in PostgreSQL.
func (r *postgresGuestRepository) AddGuest(ctx context.Context, guest *types.GuestType) (*types.GuestType, error) {
	guest.ID = uuid.New()
	if err := r.db.WithContext(ctx).Create(guest).Error; err != nil {
		return nil, err
	}
	return guest, nil
}

// AddBulkGuest adds multiple guests to an event in PostgreSQL.
func (r *postgresGuestRepository) AddBulkGuest(ctx context.Context, guests []*types.GuestType) ([]*types.GuestType, error) {
	// Generate UUIDs for all guests
	for _, guest := range guests {
		guest.ID = uuid.New()
	}

	// Batch insert guests using GORM's Create method with batch size
	batchSize := len(guests)
	if err := r.db.WithContext(ctx).CreateInBatches(guests, batchSize).Error; err != nil {
		return nil, err
	}

	return guests, nil
}

// FindGuestsByEventID retrieves all guests for a given event in PostgreSQL.
func (r *postgresGuestRepository) FindGuestsByEventID(ctx context.Context, eventID uuid.UUID) ([]*types.GuestType, error) {
	var guests []*types.GuestType
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&guests).Error; err != nil {
		return nil, err
	}
	return guests, nil
}

// FindGuestByID retrieves a guest by their ID in PostgreSQL.
func (r *postgresGuestRepository) FindGuestByID(ctx context.Context, guestID uuid.UUID) (*types.GuestType, error) {
	var guest types.GuestType
	if err := r.db.WithContext(ctx).Where("id = ?", guestID).First(&guest).Error; err != nil {
		return nil, err
	}
	return &guest, nil
}

// UpdateGuest updates the information of an existing guest in PostgreSQL.
func (r *postgresGuestRepository) UpdateGuest(ctx context.Context, guest *types.GuestType) error {
	return r.db.WithContext(ctx).Save(guest).Error
}

// DeleteGuestByID removes a guest by their ID in PostgreSQL.
func (r *postgresGuestRepository) DeleteGuestByID(ctx context.Context, guestID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", guestID).Delete(&types.GuestType{}).Error
}

// CountGuestsByEventID counts the number of guests for a given event in PostgreSQL.
func (r *postgresGuestRepository) CountGuestsByEventID(ctx context.Context, eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&types.GuestType{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}
