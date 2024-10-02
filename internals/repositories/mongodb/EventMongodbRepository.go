package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoEventRepository struct {
	collection *mongo.Collection
}

// NewMongoEventRepository initializes a new instance of the event repository.
func NewMongoEventRepository(db *mongo.Database) repositories.EventRepositoryInterface {
	return &mongoEventRepository{
		collection: db.Collection("events"),
	}
}

// CreateEvent creates a new event record in MongoDB.
func (r *mongoEventRepository) CreateEvent(ctx context.Context, event *types.EventType) (*types.EventType, error) {
	event.ID = uuid.New()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// FindEventByID finds an event by its ID in MongoDB.
func (r *mongoEventRepository) FindEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error) {
	var event types.EventType
	err := r.collection.FindOne(ctx, bson.M{"id": eventID}).Decode(&event)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &event, nil
}

// FindEventsByOrganizerID retrieves events organized by a specific user.
func (r *mongoEventRepository) FindEventsByOrganizerID(ctx context.Context, organizerID uuid.UUID) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"organizer_id": organizerID})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateEvent updates an existing event in MongoDB.
func (r *mongoEventRepository) UpdateEvent(ctx context.Context, event *types.EventType) error {
	event.UpdatedAt = time.Now()
	filter := bson.M{"id": event.ID}
	update := bson.M{"$set": event}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteEventByID deletes an event by its ID in MongoDB.
func (r *mongoEventRepository) DeleteEventByID(ctx context.Context, eventID uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": eventID})
	return err
}

// FindAllEvents retrieves all events from MongoDB.
func (r *mongoEventRepository) FindAllEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// FindEventsByDateRange finds events occurring between a start and end date.
func (r *mongoEventRepository) FindEventsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{
		"start_time": bson.M{"$gte": startDate},
		"end_time":   bson.M{"$lte": endDate},
	})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// FindEventsByLocation finds events based on location.
func (r *mongoEventRepository) FindEventsByLocation(ctx context.Context, location string) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"location": location})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// FindEventsByMultipleTags retrieves events that match any of the specified tags.
func (r *mongoEventRepository) FindEventsByMultipleTags(ctx context.Context, tags []string) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"tags": bson.M{"$in": tags}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// ActivateEvent activates an event, making it live in MongoDB.
func (r *mongoEventRepository) ActivateEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.updateEventStatus(ctx, eventID, true)
}

// DeactivateEvent deactivates an event, making it inactive in MongoDB.
func (r *mongoEventRepository) DeactivateEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.updateEventStatus(ctx, eventID, false)
}

// FindUpcomingEvents retrieves events scheduled for future dates.
func (r *mongoEventRepository) FindUpcomingEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"start_time": bson.M{"$gt": time.Now()}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// FindPastEvents retrieves events that have already occurred.
func (r *mongoEventRepository) FindPastEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"end_time": bson.M{"$lt": time.Now()}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// FindActiveEvents retrieves all currently active events.
func (r *mongoEventRepository) FindActiveEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// FindInactiveEvents retrieves all inactive events.
func (r *mongoEventRepository) FindInactiveEvents(ctx context.Context) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": false})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// SearchEventsByTitle finds events by searching their titles in MongoDB.
func (r *mongoEventRepository) SearchEventsByTitle(ctx context.Context, title string) ([]*types.EventType, error) {
	var events []*types.EventType
	cursor, err := r.collection.Find(ctx, bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// CancelEvent cancels an event by marking it as inactive in MongoDB.
func (r *mongoEventRepository) CancelEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.updateEventStatus(ctx, eventID, false)
}

// RescheduleEvent reschedules an event to a new date and time in MongoDB.
func (r *mongoEventRepository) RescheduleEvent(ctx context.Context, eventID uuid.UUID, newStartTime, newEndTime time.Time) error {
	filter := bson.M{"id": eventID}
	update := bson.M{
		"$set": bson.M{
			"start_time": newStartTime,
			"end_time":   newEndTime,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// CountTotalEvents returns the total number of events in the system.
func (r *mongoEventRepository) CountTotalEvents(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountEventsByOrganizerID counts the number of events organized by a specific user.
func (r *mongoEventRepository) CountEventsByOrganizerID(ctx context.Context, organizerID uuid.UUID) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"organizer_id": organizerID})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// updateEventStatus is a helper function to activate or deactivate an event.
func (r *mongoEventRepository) updateEventStatus(ctx context.Context, eventID uuid.UUID, isActive bool) error {
	filter := bson.M{"id": eventID}
	update := bson.M{
		"$set": bson.M{
			"is_active":  isActive,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
