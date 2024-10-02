package handlers

import "github.com/lordofthemind/EventureGo/internals/services"

type EventGinHandler struct {
	service services.EventServiceInterface
}

func NewEventGinHandler(service services.EventServiceInterface) *EventGinHandler {
	return &EventGinHandler{
		service: service,
	}
}
