package routes

import (
	"github.com/Clementol/Wallet/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.RouterGroup) {

	incomingRoutes.POST("/signup", controllers.Signup())
	incomingRoutes.POST("/signin", controllers.Login())

}
