package types

import (
	"time"

	"github.com/google/uuid"
)

// BaseUserType defines the structure for shared user fields (SuperUser, Guest)
type BaseUserType struct {
	ID        uuid.UUID `bson:"_id,omitempty" json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FullName  string    `bson:"full_name" json:"full_name" validate:"required,min=3,max=32" gorm:"not null"`
	Email     string    `bson:"email" json:"email" validate:"required,email" gorm:"unique;not null"`
	CreatedAt time.Time `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at" gorm:"autoUpdateTime"`
	IsActive  bool      `bson:"is_active" json:"is_active" gorm:"default:true"`
}
