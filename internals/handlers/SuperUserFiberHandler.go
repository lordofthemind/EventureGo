package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
)

type SuperUserFiberHandler struct {
	service services.SuperUserServiceInterface
}

func NewSuperUserFiberHandler(service services.SuperUserServiceInterface) *SuperUserFiberHandler {
	return &SuperUserFiberHandler{service: service}
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
		}
		response := responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to register superuser", nil, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Use standardized response for successful registration
	response := responses.NewFiberResponse(c, fiber.StatusCreated, "Superuser registered successfully", registeredSuperUser, nil)
	return c.Status(fiber.StatusCreated).JSON(response)
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
		}
		response := responses.NewFiberResponse(c, fiber.StatusInternalServerError, "Failed to log in", nil, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Use standardized response for successful login
	response := responses.NewFiberResponse(c, fiber.StatusOK, "Login successful", loggedInSuperUser, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
