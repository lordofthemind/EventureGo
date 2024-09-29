package postgresdb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"gorm.io/gorm"
)

type postgresSuperUserRepository struct {
	db *gorm.DB
}

// NewPostgresSuperUserRepository initializes a new instance of the repository.
func NewPostgresSuperUserRepository(db *gorm.DB) repositories.SuperUserRepositoryInterface {
	return &postgresSuperUserRepository{
		db: db,
	}
}

// CreateSuperUser creates a new superuser record in PostgreSQL.
func (r *postgresSuperUserRepository) CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error) {
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(superUser).Error; err != nil {
		return nil, err
	}
	return superUser, nil
}

// FindSuperUserByEmail searches for a superuser by email.
func (r *postgresSuperUserRepository) FindSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&superUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindSuperUserByUsername searches for a superuser by username.
func (r *postgresSuperUserRepository) FindSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&superUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindSuperUserByResetToken searches for a superuser by reset token.
func (r *postgresSuperUserRepository) FindSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.db.WithContext(ctx).Where("reset_token = ?", token).First(&superUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

func (r *postgresSuperUserRepository) UpdateResetToken(ctx context.Context, superUserID uuid.UUID, resetToken string) error {
	return r.db.WithContext(ctx).Model(&types.SuperUserType{}).Where("id = ?", superUserID).
		Updates(map[string]interface{}{
			"reset_token": resetToken,
			"updated_at":  time.Now(),
		}).Error
}

func (r *postgresSuperUserRepository) UpdateSuperUser(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(superUser).Error
}

// FindSuperUserByOTP retrieves the superuser by OTP
func (r *postgresSuperUserRepository) FindSuperUserByOTP(ctx context.Context, otp string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	if err := r.db.WithContext(ctx).Where("otp = ?", otp).First(&superUser).Error; err != nil {
		return nil, err
	}

	// Check if OTP has expired
	if time.Now().After(superUser.OTPExpiry) {
		return nil, errors.New("OTP has expired")
	}

	return &superUser, nil
}

// VerifySuperUserOTP marks the user as verified and updates the OTP status
func (r *postgresSuperUserRepository) VerifySuperUserOTP(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.IsOTPVerified = true
	superUser.OTP = nil               // Clear the OTP after verification
	superUser.OTPExpiry = time.Time{} // Reset OTP expiry
	superUser.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Save(superUser).Error
}
