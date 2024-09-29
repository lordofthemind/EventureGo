package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/middlewares"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func SetupSuperUserGinRoutes(
	router *gin.Engine,
	superUserHandler *handlers.SuperUserGinHandler,
	tokenManager gophertoken.TokenManager,
) {
	// Public routes for registration and password resets
	router.POST("/superusers/register", superUserHandler.RegisterSuperUserHandler)
	router.GET("/superusers/verify", superUserHandler.VerifySuperUserHandler)
	router.POST("/superusers/login", superUserHandler.LogInSuperUserHandler)

	// Auth routes for protected actions
	authRoutes := router.Group("/superusers")
	authRoutes.Use(middlewares.AuthTokenGinMiddleware(tokenManager)) // Middleware to protect routes
	{
		authRoutes.GET("/logout", superUserHandler.LogOutSuperUserHandler)
		// Add more authenticated routes here as needed
	}
	// Password reset routes
	router.POST("/superuser/password-reset/request", superUserHandler.PasswordResetRequestHandler)
	router.POST("/superuser/password-reset/:token", superUserHandler.PasswordResetHandler)
}
