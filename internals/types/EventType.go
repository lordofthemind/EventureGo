package types

import (
	"time"

	"github.com/google/uuid"
)

// EventType defines the structure for an event
type EventType struct {
	ID          uuid.UUID   `bson:"_id,omitempty" json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string      `bson:"title" json:"title" validate:"required,min=3,max=100" gorm:"not null"`
	Description string      `bson:"description" json:"description" gorm:"type:text"`
	StartTime   time.Time   `bson:"start_time" json:"start_time" gorm:"not null"`
	EndTime     time.Time   `bson:"end_time" json:"end_time" gorm:"not null"`
	Location    string      `bson:"location" json:"location" gorm:"not null"`
	OrganizerID uuid.UUID   `bson:"organizer_id" json:"organizer_id" gorm:"type:uuid;not null"`
	Guests      []GuestType `bson:"guests" json:"guests" gorm:"foreignKey:EventID"`
	CreatedAt   time.Time   `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at" gorm:"autoUpdateTime"`
	IsActive    bool        `bson:"is_active" json:"is_active" gorm:"default:true"`
	Tags        []string    `bson:"tags" json:"tags" gorm:"type:text[]"`
}

// NewEvent creates a new instance of EventType
func NewEvent(title, description, location string, startTime, endTime time.Time, organizerID uuid.UUID, tags []string) *EventType {
	return &EventType{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    location,
		OrganizerID: organizerID,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
	}
}
