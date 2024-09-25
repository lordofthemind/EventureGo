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
	// Bind and validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		return
	}

	// Call the service to register the superuser
	registeredSuperUser, err := h.service.RegisterSuperUser(c.Request.Context(), &req)
	if err != nil {
		// Handle validation errors returned by service
		if err.Error() == "email already in use" || err.Error() == "username already in use" {
			responses.NewGinResponse(c, http.StatusConflict, err.Error(), nil, nil)
		} else {
			responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to register superuser", nil, err.Error())
		}
		return
	}

	// Use standardized response for successful registration
	responses.NewGinResponse(c, http.StatusCreated, "Superuser registered successfully", registeredSuperUser, nil)
}

// LogInSuperUserHandler handles the login of a superuser
func (h *SuperUserGinHandler) LogInSuperUserHandler(c *gin.Context) {
	var req utils.LogInSuperuserRequest
	// Bind and validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		return
	}

	// Ensure that either email or username is provided
	if req.Email == "" && req.Username == "" {
		responses.NewGinResponse(c, http.StatusBadRequest, "Either email or username is required", nil, nil)
		return
	}

	// Call the service to log in the superuser
	loggedInSuperUser, err := h.service.LogInSuperuser(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email/username or password" {
			responses.NewGinResponse(c, http.StatusUnauthorized, "Invalid email/username or password", nil, nil)
		} else {
			responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to log in", nil, err.Error())
		}
		return
	}

	// Use standardized response for successful login
	responses.NewGinResponse(c, http.StatusOK, "Login successful", loggedInSuperUser, nil)
}
