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

func (e *EventService) CreateEventService(ctx context.Context, organizerID uuid.UUID, eventReq *utils.RegisterEventRequest) (*types.EventType, error) {
	// Create new event using the RegisterEventRequest and organizerID
	event := types.NewEvent(
		eventReq.Title,
		eventReq.Description,
		eventReq.Location,
		eventReq.StartTime,
		eventReq.EndTime,
		organizerID,
		eventReq.Tags,
	)

	// Call repository to save event to the database
	createdEvent, err := e.repository.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}
