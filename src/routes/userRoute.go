package routes

import (
	controller "github.com/AnthuanGarcia/appWebSeguridad/src/controllers"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("/signup", controller.SignUp)
	incomingRoutes.POST("/login", controller.Login)

}
