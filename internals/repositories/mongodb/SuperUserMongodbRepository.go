package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoSuperUserRepository struct {
	collection *mongo.Collection
}

func NewMongoSuperUserRepository(db *mongo.Database) repositories.SuperUserRepositoryInterface {
	return &mongoSuperUserRepository{
		collection: db.Collection("superusers"),
	}
}

func (r *mongoSuperUserRepository) CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error) {
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, superUser)
	return superUser, err
}

func (r *mongoSuperUserRepository) FindSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

func (r *mongoSuperUserRepository) FindSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindSuperuserByResetToken finds a superuser by reset token.
func (r *mongoSuperUserRepository) FindSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superuser types.SuperUserType
	err := r.collection.FindOne(ctx, bson.M{"reset_token": token}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}
