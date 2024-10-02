package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"github.com/lordofthemind/EventureGo/internals/utils"
)

type EventService struct {
	repository repositories.EventRepositoryInterface
}

func NewEventService(repository repositories.EventRepositoryInterface) EventServiceInterface {
	return &EventService{
		repository: repository,
	}
}

func (e *EventService) CreateEventService(ctx context.Context, organizerID uuid.UUID, eventDTO *utils.EventDTO) (*types.EventType, error) {
	// Create new event using the DTO and organizerID
	event := types.NewEvent(
		eventDTO.Title,
		eventDTO.Description,
		eventDTO.Location,
		eventDTO.StartTime,
		eventDTO.EndTime,
		organizerID,
		eventDTO.Tags,
	)

	// Call repository to save the event in the database
	createdEvent, err := e.repository.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}
