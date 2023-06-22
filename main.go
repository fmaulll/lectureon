package main

import (
	"net/http"
	"os"

	"github.com/fmaulll/lectureon/controllers"
	"github.com/fmaulll/lectureon/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
	initializers.Migrate()
}

func main() {
	router := gin.Default()

	router.GET("/api/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "WOW"})
	})

	router.POST("/api/signup", controllers.Signup)

	router.Run(":" + os.Getenv("PORT"))
}