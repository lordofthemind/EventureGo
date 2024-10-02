package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/configs"
)

// SuperUserType extends BaseUserType for superusers
type SuperUserType struct {
	BaseUserType
	Role             string    `bson:"role" json:"role" validate:"required" gorm:"not null"`
	Username         string    `bson:"username" json:"username" validate:"required,min=3,max=32,alphanum" gorm:"unique;not null"`
	HashedPassword   string    `bson:"hashed_password" json:"-" validate:"required,min=8" gorm:"not null"`
	ResetToken       *string   `bson:"reset_token,omitempty" json:"reset_token,omitempty" gorm:"type:text"`
	ResetTokenExpiry time.Time `bson:"reset_token_expiry,omitempty" json:"-"`
	Is2FAEnabled     bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled" gorm:"default:false"`
	TwoFactorSecret  *string   `bson:"two_factor_secret,omitempty" json:"-" gorm:"type:text"`
	PermissionGroups []string  `bson:"permission_groups" json:"permission_groups" validate:"dive,required" gorm:"type:text[]"`
	OTP              *string   `bson:"otp,omitempty" json:"-" gorm:"type:text"`
	OTPExpiry        time.Time `bson:"otp_expiry,omitempty" json:"-"`
	IsOTPVerified    bool      `bson:"is_otp_verified" json:"is_otp_verified" gorm:"default:false"`
}

// NewSuperUser creates a new SuperUser instance
func NewSuperUser(email, fullName, username, hashedPassword, role, otp string) *SuperUserType {
	return &SuperUserType{
		BaseUserType: BaseUserType{
			ID:        uuid.New(),
			Email:     email,
			FullName:  fullName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsActive:  true,
		},
		Role:           validateRole(role),
		Username:       username,
		HashedPassword: hashedPassword,
		OTP:            &otp,
		OTPExpiry:      time.Now().Add(15 * time.Minute), // Example OTP expiry
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
