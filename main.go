package main

import (
	"example/REST_API_JWT/controllers"
	"example/REST_API_JWT/db"
	"example/REST_API_JWT/initializers"
	"example/REST_API_JWT/middlewares"

	// "fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	db.ConnectToDb()

	// sync the database
	db.SyncDatabase()
}

func main() {
	router := gin.Default()
	// test-api
	router.GET("/ping", func(res *gin.Context) {
		res.JSON(200, gin.H{
			"message": "Pinged Successfully",
			"success": true,
		})
	})

	router.POST("/api/v1/auth/create-user", controllers.Signup)

	router.POST("/api/v1/auth/login", controllers.SingIn)

	router.GET("/api/v1/auth/get-user-details", middlewares.DecodeJwt, controllers.GetUserDetails)
	router.Run()
}
