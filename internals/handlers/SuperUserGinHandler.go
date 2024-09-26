package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
	"github.com/lordofthemind/mygopher/gophertoken"
)

type SuperUserGinHandler struct {
	service      services.SuperUserServiceInterface
	tokenManager gophertoken.TokenManager
}

func NewSuperUserGinHandler(service services.SuperUserServiceInterface, tokenManager gophertoken.TokenManager) *SuperUserGinHandler {
	return &SuperUserGinHandler{
		service:      service,
		tokenManager: tokenManager,
	}
}

// RegisterSuperUserHandler handles the registration of a new superuser
func (h *SuperUserGinHandler) RegisterSuperUserHandler(c *gin.Context) {
	var req utils.RegisterSuperuserRequest
	// Bind and validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Call the service to register the superuser
	registeredSuperUser, err := h.service.RegisterSuperUser(c.Request.Context(), &req)
	if err != nil {
		// Handle validation errors returned by service
		if err.Error() == "email already in use" || err.Error() == "username already in use" {
			response := responses.NewGinResponse(c, http.StatusConflict, err.Error(), nil, nil)
			c.JSON(http.StatusConflict, response)
		} else {
			response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to register superuser", nil, err.Error())
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	// Use standardized response for successful registration
	response := responses.NewGinResponse(c, http.StatusCreated, "Superuser registered successfully", registeredSuperUser, nil)
	c.JSON(http.StatusCreated, response)
}

// LogInSuperUserHandler handles the login of a superuser
func (h *SuperUserGinHandler) LogInSuperUserHandler(c *gin.Context) {
	var req utils.LogInSuperuserRequest
	// Bind and validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Ensure that either email or username is provided
	if req.Email == "" && req.Username == "" {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Either email or username is required", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Call the service to log in the superuser
	loggedInSuperUser, err := h.service.LogInSuperuser(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email/username or password" {
			response := responses.NewGinResponse(c, http.StatusUnauthorized, "Invalid email/username or password", nil, nil)
			c.JSON(http.StatusUnauthorized, response)
		} else {
			response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to log in", nil, err.Error())
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	authToken, err := h.tokenManager.GenerateToken(loggedInSuperUser.Username, configs.TokenExpiryDuration)
	if err != nil {
		responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to generate token", nil, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("SuperUserAuthorizationToken", authToken, int(configs.TokenExpiryDuration.Seconds()), "/", "", false, true)

	// Use standardized response for successful login
	response := responses.NewGinResponse(c, http.StatusOK, "Login successful", loggedInSuperUser, nil)
	c.JSON(http.StatusOK, response)
}

func (h *SuperUserGinHandler) LogOutSuperUserHandler(c *gin.Context) {
	c.SetCookie("SuperUserAuthorizationToken", "", -1, "/", "", false, true)
	// Use standardized response for successful login
	response := responses.NewGinResponse(c, http.StatusOK, "Logout successful", nil, nil)
	c.JSON(http.StatusOK, response)
}
