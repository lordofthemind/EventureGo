package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/middlewares"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func SetupSuperUserGinRoutes(router *gin.Engine, superUserHandler *handlers.SuperUserGinHandler, tokenManager gophertoken.TokenManager) {
	// Public routes
	router.POST("/superusers/register", superUserHandler.RegisterSuperUserHandler)
	router.POST("/superusers/login", superUserHandler.LogInSuperUserHandler)

	// Protected routes with JWTAuthMiddleware
	protectedRoutes := router.Group("/superusers")
	protectedRoutes.Use(middlewares.AuthTokenMiddleware(tokenManager))
	{
		protectedRoutes.GET("/logout", superUserHandler.LogOutSuperUserHandler)
	}
}
