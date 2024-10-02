package validators

import (
	"errors"
	"time"

	"github.com/lordofthemind/EventureGo/internals/utils"
)

// ValidateEventRequest checks if the event data is valid
func ValidateEventRequest(req utils.RegisterEventRequest) error {
	if len(req.Title) < 3 || len(req.Title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}
	if req.StartTime.After(req.EndTime) {
		return errors.New("start time cannot be after end time")
	}
	if time.Now().After(req.StartTime) {
		return errors.New("event start time cannot be in the past")
	}
	return nil
}
