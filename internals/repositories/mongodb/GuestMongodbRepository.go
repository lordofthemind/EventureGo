package mongodb

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoGuestRepository struct {
	collection *mongo.Collection
}

// NewMongoGuestRepository initializes a new instance of the guest repository.
func NewMongoGuestRepository(db *mongo.Database) repositories.GuestRepositoryInterface {
	return &mongoGuestRepository{
		collection: db.Collection("guests"),
	}
}

// AddGuest adds a new guest to an event in MongoDB.
func (r *mongoGuestRepository) AddGuest(ctx context.Context, guest *types.GuestType) (*types.GuestType, error) {
	guest.ID = uuid.New()
	_, err := r.collection.InsertOne(ctx, guest)
	if err != nil {
		return nil, err
	}
	return guest, nil
}

// AddBulkGuest adds multiple guests to an event in MongoDB.
func (r *mongoGuestRepository) AddBulkGuest(ctx context.Context, guests []*types.GuestType) ([]*types.GuestType, error) {
	var bulkOps []mongo.WriteModel
	for _, guest := range guests {
		guest.ID = uuid.New()
		bulkOps = append(bulkOps, mongo.NewInsertOneModel().SetDocument(guest))
	}
	_, err := r.collection.BulkWrite(ctx, bulkOps)
	if err != nil {
		return nil, err
	}
	return guests, nil
}

// FindGuestsByEventID retrieves all guests for a given event in MongoDB.
func (r *mongoGuestRepository) FindGuestsByEventID(ctx context.Context, eventID uuid.UUID) ([]*types.GuestType, error) {
	filter := bson.M{"event_id": eventID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var guests []*types.GuestType
	if err = cursor.All(ctx, &guests); err != nil {
		return nil, err
	}
	return guests, nil
}

// FindGuestByID retrieves a guest by their ID in MongoDB.
func (r *mongoGuestRepository) FindGuestByID(ctx context.Context, guestID uuid.UUID) (*types.GuestType, error) {
	filter := bson.M{"id": guestID}
	var guest types.GuestType
	err := r.collection.FindOne(ctx, filter).Decode(&guest)
	if err != nil {
		return nil, err
	}
	return &guest, nil
}

// UpdateGuest updates the information of an existing guest in MongoDB.
func (r *mongoGuestRepository) UpdateGuest(ctx context.Context, guest *types.GuestType) error {
	filter := bson.M{"id": guest.ID}
	update := bson.M{"$set": guest}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteGuestByID removes a guest by their ID in MongoDB.
func (r *mongoGuestRepository) DeleteGuestByID(ctx context.Context, guestID uuid.UUID) error {
	filter := bson.M{"id": guestID}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// CountGuestsByEventID counts the number of guests for a given event in MongoDB.
func (r *mongoGuestRepository) CountGuestsByEventID(ctx context.Context, eventID uuid.UUID) (int64, error) {
	filter := bson.M{"event_id": eventID}
	count, err := r.collection.CountDocuments(ctx, filter)
	return count, err
}
