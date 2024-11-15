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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	if err != nil {
		return nil, err
	}
	return superUser, nil
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

func (r *mongoSuperUserRepository) FindSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.collection.FindOne(ctx, bson.M{"reset_token": token}).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

func (r *mongoSuperUserRepository) UpdateResetToken(ctx context.Context, superUserID uuid.UUID, resetToken string, resetTokenExpiry time.Time) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": superUserID},
		bson.M{
			"$set": bson.M{
				"reset_token":        resetToken,
				"reset_token_expiry": resetTokenExpiry,
				"updated_at":         time.Now(),
			},
		},
	)
	return err
}

func (r *mongoSuperUserRepository) UpdateSuperUser(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"id": superUser.ID},
		superUser,
		options.Replace().SetUpsert(true),
	)
	return err
}

// FindSuperUserByOTP retrieves the superuser by OTP from MongoDB
func (r *mongoSuperUserRepository) FindSuperUserByOTP(ctx context.Context, otp string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"otp": otp}
	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err != nil {
		return nil, err
	}

	// Check if OTP has expired
	if time.Now().After(superUser.OTPExpiry) {
		return nil, errors.New("OTP has expired")
	}

	return &superUser, nil
}

// VerifySuperUserOTP marks the user as verified in MongoDB
func (r *mongoSuperUserRepository) VerifySuperUserOTP(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.IsOTPVerified = true
	superUser.OTP = nil               // Clear the OTP after verification
	superUser.OTPExpiry = time.Time{} // Reset OTP expiry
	superUser.UpdatedAt = time.Now()

	filter := bson.M{"_id": superUser.ID}
	update := bson.M{
		"$set": bson.M{
			"is_otp_verified": true,
			"otp":             nil,
			"otp_expiry":      time.Time{},
			"updated_at":      superUser.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
