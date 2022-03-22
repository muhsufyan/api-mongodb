package main

import (
	"os"

	"github.com/gin-gonic/gin"
	routes "github.com/muhsufyan/api-mongodb/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	// middleware bawaan logger
	router.Use(gin.Logger())

	// middleware yg kita buat di routes/*
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-v1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "masuk",
		})
	})
	router.GET("/api-v2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "masuk",
		})
	})
	router.Run(":" + port)
}
