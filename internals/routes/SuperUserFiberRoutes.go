package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/middlewares"
	"github.com/lordofthemind/mygopher/gophertoken"
)

// SetupSuperUserFiberRoutes sets up routes for the superuser operations in Fiber
func SetupSuperUserFiberRoutes(
	app *fiber.App,
	handler *handlers.SuperUserFiberHandler,
	tokenManager gophertoken.TokenManager,
) {
	// Public routes for registration and login
	app.Post("/superusers/register", handler.RegisterSuperUserHandler)
	app.Get("/superusers/verify", handler.VerifySuperUserHandler)
	app.Post("/superusers/login", handler.LogInSuperUserHandler)

	// Authenticated routes group, protected by the AuthTokenFiberMiddleware
	authRoutes := app.Group("/superusers", middlewares.AuthTokenFiberMiddleware(tokenManager))

	// Authenticated actions (e.g., logout)
	authRoutes.Get("/logout", handler.LogOutSuperUserHandler)
	// Add more authenticated routes here as needed

	// Password reset routes
	app.Post("/superuser/password-reset/request", handler.PasswordResetRequestHandler)
	app.Post("/superuser/password-reset/:token", handler.PasswordResetHandler)
}
