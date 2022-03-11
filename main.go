package main

import (
	"os"

	"github.com/Clementol/Wallet/routes"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	version1 := router.Group("api/v1")
	// router.Static("/public", "./public/images")

	routes.UserRoutes(version1)
	routes.WalletRoutes(version1)

	router.SetTrustedProxies(nil)

	router.Run(":" + PORT)
}
