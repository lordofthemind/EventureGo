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

	// Generate dynamic cookie name using the superuser's role
	cookieName := loggedInSuperUser.Role + "|_|" + configs.TokenBaseCookieName

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		cookieName,
		loggedInSuperUser.Token,
		int(configs.TokenExpiryDuration.Seconds()),
		"/",
		"",
		configs.SecureCookieHTTPS,
		true,
	)

	// Use standardized response for successful login
	response := responses.NewGinResponse(c, http.StatusOK, "Login successful", loggedInSuperUser, nil)
	c.JSON(http.StatusOK, response)
}

func (h *SuperUserGinHandler) LogOutSuperUserHandler(c *gin.Context) {
	// Extract the role from the context
	role, exists := c.Get("role")
	if !exists {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Role not found in context", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create the cookie name using the role and the base cookie name
	cookieName := role.(string) + "|_|" + configs.TokenBaseCookieName

	// Clear the cookie by setting its expiration to a past time
	c.SetCookie(cookieName, "", -1, "/", "", false, true)

	// Use standardized response for successful logout
	response := responses.NewGinResponse(c, http.StatusOK, "Logout successful", nil, nil)
	c.JSON(http.StatusOK, response)
}

// PasswordResetRequestHandler handles the request to send a password reset email
func (h *SuperUserGinHandler) PasswordResetRequestHandler(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"omitempty,email" validate:"omitempty,email"`
		Username string `json:"username" binding:"omitempty,min=3" validate:"omitempty,min=3"`
	}

	// Bind and validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid email or username", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Ensure that either email or username is provided
	if request.Email == "" && request.Username == "" {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Either email or username is required", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Call the service to send a password reset email or username reset token
	err := h.service.SendPasswordResetEmailWithUsernameOrEmail(c.Request.Context(), request.Email, request.Username)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to send reset email", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Use standardized response for successful email request
	response := responses.NewGinResponse(c, http.StatusOK, "Password reset email sent successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

// PasswordResetHandler handles the password reset using a token
func (h *SuperUserGinHandler) PasswordResetHandler(c *gin.Context) {
	var request struct {
		Password string `json:"password" binding:"required,min=8"`
	}

	// Bind and validate the request payload
	if err := c.ShouldBindJSON(&request); err != nil {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid password", nil, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Extract the reset token from the URL parameters
	token := c.Param("token")
	if token == "" {
		response := responses.NewGinResponse(c, http.StatusBadRequest, "Invalid reset token", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Call the service to reset the password
	err := h.service.ResetPassword(c.Request.Context(), token, request.Password)
	if err != nil {
		response := responses.NewGinResponse(c, http.StatusInternalServerError, "Failed to reset password", nil, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Use standardized response for successful password reset
	response := responses.NewGinResponse(c, http.StatusOK, "Password reset successful", nil, nil)
	c.JSON(http.StatusOK, response)
}
