package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	var req utils.RegisterSuperuserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Use standardized response for invalid input
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Call the service to register the superuser
	registeredSuperUser, err := h.service.RegisterSuperUser(c.Request.Context(), &req)
	if err != nil {
		// Use standardized response for internal server errors
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to register superuser", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Use standardized response for successful registration
	response := responses.NewGinResponse(c, http.StatusCreated, "Superuser registered successfully", registeredSuperUser, nil)
	c.JSON(http.StatusCreated, response)
}
