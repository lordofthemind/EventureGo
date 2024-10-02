package types

import (
	"time"

	"github.com/google/uuid"
)

// GuestType extends BaseUserType for guests
type GuestType struct {
	BaseUserType
	RSVPStatus string    `bson:"rsvp_status" json:"rsvp_status" gorm:"type:text"`
	EventID    uuid.UUID `bson:"event_id" json:"event_id" gorm:"not null"`
	InvitedAt  time.Time `bson:"invited_at" json:"invited_at" gorm:"autoCreateTime"`
}

// NewGuest creates a new Guest instance
func NewGuest(email, fullName string) *GuestType {
	return &GuestType{
		BaseUserType: BaseUserType{
			ID:        uuid.New(),
			Email:     email,
			FullName:  fullName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsActive:  true,
		},
		RSVPStatus: "Pending",
	}
}
