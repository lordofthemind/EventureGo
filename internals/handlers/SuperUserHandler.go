package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/internals/newerrors"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
)

type SuperUserGinHandler struct {
	service services.SuperUserServiceInterface
}

func NewSuperUserGinHandler(service services.SuperUserServiceInterface) *SuperUserGinHandler {
	return &SuperUserGinHandler{service: service}
}

// RegisterSuperUserHandler handles the registration of a new superuser
func (h *SuperUserGinHandler) RegisterSuperUserHandler(c *gin.Context) {
	// Assuming validation middleware is applied to validate the request
	var req utils.RegisterSuperuserRequest
	if err := c.ShouldBind(&req); err != nil {
		responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		return
	}

	// Call the service to register the superuser
	registeredSuperUser, err := h.service.RegisterSuperUser(c.Request.Context(), &req)
	if err != nil {
		if newerrors.IsValidationError(err) {
			responses.NewGinResponse(c, http.StatusConflict, err.Error(), nil, nil)
		} else {
			responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to register superuser", nil, err.Error())
		}
		return
	}

	// Successful registration
	responses.NewGinResponse(c, http.StatusCreated, "Superuser registered successfully", registeredSuperUser, nil)
}
