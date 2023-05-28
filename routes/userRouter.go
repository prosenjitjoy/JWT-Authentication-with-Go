package routes

import (
	"github.com/gin-gonic/gin"
	"main/controllers"
	"main/middleware"
)

func UserRoutes(inRoutes *gin.Engine) {
	inRoutes.Use(middleware.Authenticate())
	inRoutes.GET("/users", controller.GetUsers())
	inRoutes.GET("/users/:user_id", controller.GetUser())
}
