package utils

import (
	"time"

	"github.com/google/uuid"
)

type RegisterSuperuserRequest struct {
	Email    string `bson:"email" json:"email" validate:"required,email" gorm:"unique;not null"`
	FullName string `bson:"full_name" json:"full_name" validate:"required,min=3,max=32" gorm:"not null"`
	Username string `bson:"username" json:"username" validate:"required,min=3,max=32,alphanum" gorm:"unique;not null"`
	Password string `bson:"password" json:"password" validate:"required,min=8" gorm:"not null"`
}

type SuperuserResponse struct {
	ID           uuid.UUID `bson:"_id,omitempty" json:"id,omitempty"`
	Role         string    `bson:"role" json:"role"`
	Email        string    `bson:"email" json:"email"`
	FullName     string    `bson:"full_name" json:"full_name"`
	Username     string    `bson:"username" json:"username"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
	Is2FAEnabled bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled"`
}
