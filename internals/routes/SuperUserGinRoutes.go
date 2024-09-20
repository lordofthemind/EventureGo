package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/internals/handlers"
)

func SetupSuperUserGinRoutes(r *gin.Engine, handler *handlers.SuperUserGinHandler) {
	r.POST("/superusers/register", handler.RegisterSuperUserHandler)
}
