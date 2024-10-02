package postgresdb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"gorm.io/gorm"
)

type postgresEventRepository struct {
	db *gorm.DB
}

// NewPostgresEventRepository initializes a new instance of the event repository.
func NewPostgresEventRepository(db *gorm.DB) repositories.EventRepositoryInterface {
	return &postgresEventRepository{
		db: db,
	}
}

// CreateEvent creates a new event record in PostgreSQL.
func (r *postgresEventRepository) CreateEvent(ctx context.Context, event *types.EventType) (*types.EventType, error) {
	event.ID = uuid.New()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

// FindEventByID finds an event by its ID in PostgreSQL.
func (r *postgresEventRepository) FindEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error) {
	var event types.EventType
	if err := r.db.WithContext(ctx).First(&event, "id = ?", eventID).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// FindEventsByOrganizerID retrieves events organized by a specific user.
func (r *postgresEventRepository) FindEventsByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("organizer_id = ?", organizerID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateEvent updates an existing event in PostgreSQL.
func (r *postgresEventRepository) UpdateEvent(ctx context.Context, event *types.EventType) error {
	event.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		return err
	}
	return nil
}

// DeleteEventByID deletes an event by its ID in PostgreSQL.
func (r *postgresEventRepository) DeleteEventByID(ctx context.Context, eventID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&types.EventType{}, "id = ?", eventID).Error; err != nil {
		return err
	}
	return nil
}

// FindAllEvents retrieves all events from PostgreSQL.
func (r *postgresEventRepository) FindAllEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// FindEventsByDateRange finds events occurring between a start and end date.
func (r *postgresEventRepository) FindEventsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("start_time >= ? AND end_time <= ?", startDate, endDate).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// FindEventsByLocation finds events based on location.
func (r *postgresEventRepository) FindEventsByLocation(ctx context.Context, location string) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("location = ?", location).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// FindEventsByMultipleTags retrieves events that match any of the specified tags.
func (r *postgresEventRepository) FindEventsByMultipleTags(ctx context.Context, tags []string) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("tags && ?", tags).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// ActivateEvent activates an event, making it live in PostgreSQL.
func (r *postgresEventRepository) ActivateEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.updateEventStatus(ctx, eventID, true)
}

// DeactivateEvent deactivates an event, making it inactive in PostgreSQL.
func (r *postgresEventRepository) DeactivateEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.updateEventStatus(ctx, eventID, false)
}

// FindUpcomingEvents retrieves events scheduled for future dates.
func (r *postgresEventRepository) FindUpcomingEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("start_time > ?", time.Now()).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// FindPastEvents retrieves events that have already occurred.
func (r *postgresEventRepository) FindPastEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("end_time < ?", time.Now()).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// FindActiveEvents retrieves all currently active events.
func (r *postgresEventRepository) FindActiveEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("is_active = true").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// FindInactiveEvents retrieves all inactive events.
func (r *postgresEventRepository) FindInactiveEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("is_active = false").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// SearchEventsByTitle finds events by searching their titles in PostgreSQL.
func (r *postgresEventRepository) SearchEventsByTitle(ctx context.Context, title string) ([]*types.EventType, error) {
	var events []*types.EventType
	if err := r.db.WithContext(ctx).Where("title ILIKE ?", "%"+title+"%").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// CancelEvent cancels an event by marking it as inactive in PostgreSQL.
func (r *postgresEventRepository) CancelEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.updateEventStatus(ctx, eventID, false)
}

// RescheduleEvent reschedules an event to a new date and time in PostgreSQL.
func (r *postgresEventRepository) RescheduleEvent(ctx context.Context, eventID uuid.UUID, newStartTime, newEndTime time.Time) error {
	var event types.EventType
	if err := r.db.WithContext(ctx).First(&event, "id = ?", eventID).Error; err != nil {
		return err
	}
	event.StartTime = newStartTime
	event.EndTime = newEndTime
	event.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(&event).Error; err != nil {
		return err
	}
	return nil
}

// CountTotalEvents returns the total number of events in the system.
func (r *postgresEventRepository) CountTotalEvents(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&types.EventType{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountEventsByOrganizerID counts the number of events organized by a specific user.
func (r *postgresEventRepository) CountEventsByOrganizerID(ctx context.Context, organizerID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&types.EventType{}).Where("organizer_id = ?", organizerID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// updateEventStatus is a helper function to activate or deactivate an event.
func (r *postgresEventRepository) updateEventStatus(ctx context.Context, eventID uuid.UUID, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&types.EventType{}).Where("id = ?", eventID).Update("is_active", isActive).Error; err != nil {
		return err
	}
	return nil
}
