package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/middlewares"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func SetupEventGinRoutes(
	router *gin.Engine,
	eventGinHandler *handlers.EventGinHandler,
	tokenManager gophertoken.TokenManager,
) {
	// Event routes for protected actions
	protectedEventRoutes := router.Group("/event")
	protectedEventRoutes.Use(middlewares.AuthTokenGinMiddleware(tokenManager)) // Middleware to protect routes
	{
		protectedEventRoutes.POST("/register", eventGinHandler.CreateEventHandler)
	}
}
