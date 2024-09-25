package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/internals/handlers"
)

// SetupSuperUserFiberRoutes sets up routes for the superuser operations in Fiber
func SetupSuperUserFiberRoutes(app *fiber.App, handler *handlers.SuperUserFiberHandler) {
	app.Post("/superusers/register", handler.RegisterSuperUserHandler)
	app.Post("/superusers/login", handler.LogInSuperUserHandler)
}
