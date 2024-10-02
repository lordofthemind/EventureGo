package services

import "github.com/lordofthemind/EventureGo/internals/repositories"

type EventService struct {
	repository repositories.EventRepositoryInterface
}

// func NewEventService(repository repositories.EventRepositoryInterface) EventServiceInterface {
// 	return &EventService{
// 		repository: repository,
// 	}
// }
