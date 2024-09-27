package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/configs"
)

// SuperUserType defines the structure for a superuser
type SuperUserType struct {
	ID               uuid.UUID `bson:"_id,omitempty" json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Role             string    `bson:"role" json:"role" validate:"required" gorm:"not null"`
	Email            string    `bson:"email" json:"email" validate:"required,email" gorm:"unique;not null"`
	FullName         string    `bson:"full_name" json:"full_name" validate:"required,min=3,max=32" gorm:"not null"`
	Username         string    `bson:"username" json:"username" validate:"required,min=3,max=32,alphanum" gorm:"unique;not null"`
	HashedPassword   string    `bson:"hashed_password" json:"-" validate:"required,min=8" gorm:"not null"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at" gorm:"autoUpdateTime"`
	ResetToken       *string   `bson:"reset_token,omitempty" json:"reset_token,omitempty" gorm:"type:text"`
	Is2FAEnabled     bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled" gorm:"default:false"`
	TwoFactorSecret  *string   `bson:"two_factor_secret,omitempty" json:"-" gorm:"type:text"`
	PermissionGroups []string  `bson:"permission_groups" json:"permission_groups" validate:"dive,required" gorm:"type:text[]"`
}

// NewSuperUser creates a new instance of SuperUserType
func NewSuperUser(email, fullName, username, hashedPassword, role string) *SuperUserType {
	return &SuperUserType{
		ID:             uuid.New(),
		Email:          email,
		FullName:       fullName,
		Username:       username,
		HashedPassword: hashedPassword,
		Role:           validateRole(role),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// ValidateRole checks if the provided role is in the allowed roles
func validateRole(role string) string {
	for _, allowedRole := range configs.AllowedRoles {
		if role == allowedRole {
			return role
		}
	}
	return "Guest"
}
