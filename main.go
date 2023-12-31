package main

import (
	"fmt"
	"os"

	"github.com/fmaulll/lectureon/controllers"
	"github.com/fmaulll/lectureon/initializers"
	"github.com/fmaulll/lectureon/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
	initializers.Migrate()
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/images", "./images")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	router.GET("/api/", func(ctx *gin.Context) {
		fmt.Println("The URL: ", ctx.Request.Host)
	})

	router.POST("/api/token", controllers.Token)

	router.POST("/api/signup", controllers.Signup)
	router.POST("/api/login", controllers.Login)

	router.POST("/api/post", middleware.RequireAuth, controllers.NewPost)
	router.GET("/api/post", middleware.RequireAuth, controllers.GetAllPost)
	router.GET("/api/post/:id", middleware.RequireAuth, controllers.GetAllPostByAuthorId)
	router.PATCH("/api/post", middleware.RequireAuth, controllers.EditPost)

	router.Run(":" + os.Getenv("PORT"))
}
