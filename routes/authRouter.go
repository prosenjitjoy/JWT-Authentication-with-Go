package routes

import (
	"github.com/gin-gonic/gin"
	"main/controllers"
)

func AuthRoutes(inRoutes *gin.Engine) {
	inRoutes.POST("users/signup", controller.Signup())
	inRoutes.POST("users/login", controller.Login())
}
