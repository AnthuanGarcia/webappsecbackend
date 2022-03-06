package main

import (
	routes "github.com/AnthuanGarcia/appWebSeguridad/src/routes"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	//_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(cors.New(
		cors.Config{
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"GET", "POST"},
			AllowHeaders:  []string{"Content-Type"},
			ExposeHeaders: []string{"Content-Length"},
		},
	))

	routes.UserRoutes(router)

	router.Run(":8080")

}
