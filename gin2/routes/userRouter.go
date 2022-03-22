package routes

import (
	controller "github.com/muhsufyan/api-mongodb/controllers"

	"github.com/gin-gonic/gin"
	"github.com/muhsufyan/api-mongodb/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/user/:user_id", controller.GetUser())
}
