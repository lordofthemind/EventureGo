package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
)

type SuperUserFiberHandler struct {
	service services.SuperUserServiceInterface
}

func NewSuperUserFiberHandler(service services.SuperUserServiceInterface) *SuperUserFiberHandler {
	return &SuperUserFiberHandler{
		service: service,
	}
}

// RegisterSuperUserHandler handles the registration of a new superuser
func (h *SuperUserFiberHandler) RegisterSuperUserHandler(c *fiber.Ctx) error {
	var req utils.RegisterSuperuserRequest
	// Parse and validate request payload
	if err := c.BodyParser(&req); err != nil {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Call the service to register the superuser
	registeredSuperUser, err := h.service.RegisterSuperUser(c.Context(), &req)
	if err != nil {
		// Handle validation errors returned by the service
		if err.Error() == "email already in use" || err.Error() == "username already in use" {
			response := responses.NewFiberResponse(c, fiber.StatusConflict, err.Error(), nil, nil)
			return c.Status(fiber.StatusConflict).JSON(response)
		} else {
			response := responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to register superuser", nil, err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Use standardized response for successful registration
	response := responses.NewFiberResponse(c, fiber.StatusCreated, "Superuser registered successfully", registeredSuperUser, nil)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// VerifySuperUserHandler handles the OTP verification
func (h *SuperUserFiberHandler) VerifySuperUserHandler(c *fiber.Ctx) error {
	otp := c.Query("otp")
	if otp == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "OTP is required"})
	}

	superUser, err := h.service.VerifySuperUserOTP(c.Context(), otp)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Account successfully verified", "superuser": superUser})
}

// LogInSuperUserHandler handles the login of a superuser
func (h *SuperUserFiberHandler) LogInSuperUserHandler(c *fiber.Ctx) error {
	var req utils.LogInSuperuserRequest
	// Parse and validate request payload
	if err := c.BodyParser(&req); err != nil {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid input", nil, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ensure that either email or username is provided
	if req.Email == "" && req.Username == "" {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Either email or username is required", nil, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Call the service to log in the superuser
	loggedInSuperUser, err := h.service.LogInSuperuser(c.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email/username or password" {
			response := responses.NewFiberResponse(c, fiber.StatusUnauthorized, "Invalid email/username or password", nil, nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		} else {
			response := responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to log in", nil, err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Create a cookie name using the superuser's role
	cookieName := loggedInSuperUser.Role + "|_|" + configs.TokenBaseCookieName

	// Create a cookie and set it in the response
	cookie := new(fiber.Cookie)
	cookie.Name = cookieName
	cookie.Value = loggedInSuperUser.Token
	cookie.Expires = time.Now().Add(configs.TokenExpiryDuration)
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	cookie.Secure = configs.SecureCookieHTTPS

	c.Cookie(cookie)

	// Use standardized response for successful login
	response := responses.NewFiberResponse(c, fiber.StatusOK, "Login successful", loggedInSuperUser, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// LogOutSuperUserHandler handles the logout of a superuser
func (h *SuperUserFiberHandler) LogOutSuperUserHandler(c *fiber.Ctx) error {
	// Extract the role from the context
	role, exists := c.Locals("role").(string)
	if !exists {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Role not found in context", nil, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create the cookie name using the role and the base cookie name
	cookieName := role + "|_|" + configs.TokenBaseCookieName

	// Clear the authorization cookie
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to a past time to clear the cookie
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
	})

	// Use standardized response for successful logout
	response := responses.NewFiberResponse(c, fiber.StatusOK, "Logout successful", nil, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// PasswordResetRequestHandler handles the request to send a password reset email
func (h *SuperUserFiberHandler) PasswordResetRequestHandler(c *fiber.Ctx) error {
	var request struct {
		Email    string `json:"email" binding:"omitempty,email" validate:"omitempty,email"`
		Username string `json:"username" binding:"omitempty,min=3" validate:"omitempty,min=3"`
	}

	// Parse and validate the request payload
	if err := c.BodyParser(&request); err != nil {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid email or username", nil, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ensure that either email or username is provided
	if request.Email == "" && request.Username == "" {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Either email or username is required", nil, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Call the service to send a password reset email
	err := h.service.SendPasswordResetEmailWithUsernameOrEmail(c.Context(), request.Email, request.Username)
	if err != nil {
		response := responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to send reset email", nil, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Use standardized response for successful email request
	response := responses.NewFiberResponse(c, fiber.StatusOK, "Password reset email sent successfully", nil, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// PasswordResetHandler handles the password reset using a token
func (h *SuperUserFiberHandler) PasswordResetHandler(c *fiber.Ctx) error {
	var request struct {
		Password string `json:"password" binding:"required,min=8"`
	}

	// Parse and validate the request payload
	if err := c.BodyParser(&request); err != nil {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid password", nil, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Extract the reset token from the URL parameters
	token := c.Params("token")
	if token == "" {
		response := responses.NewFiberResponse(c, fiber.StatusBadRequest, "Invalid reset token", nil, nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Call the service to reset the password
	err := h.service.ResetPassword(c.Context(), token, request.Password)
	if err != nil {
		response := responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to reset password", nil, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Use standardized response for successful password reset
	response := responses.NewFiberResponse(c, fiber.StatusOK, "Password reset successful", nil, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
