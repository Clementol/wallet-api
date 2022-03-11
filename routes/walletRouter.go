package routes

import (
	"github.com/Clementol/Wallet/controllers"
	"github.com/Clementol/Wallet/middleware"
	"github.com/gin-gonic/gin"
)

func WalletRoutes(incomingRoutes *gin.RouterGroup) {

	incomingRoutes.PUT("/fund", middleware.Authentication(), controllers.FundWallet())
	// incomingRoutes.POST("/signin", controllers.Login())

}
