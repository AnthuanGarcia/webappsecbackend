package routes

import (
	controller "github.com/AnthuanGarcia/appWebSeguridad/src/controllers"
	middleware "github.com/AnthuanGarcia/appWebSeguridad/src/middleware"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("/signup", controller.SignUp)
	incomingRoutes.POST("/login", controller.Login)
	incomingRoutes.POST("/many", controller.ManyUsers)
	incomingRoutes.GET("/dash", middleware.Authentication, controller.Dashboard)

}
