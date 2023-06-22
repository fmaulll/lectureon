package controllers

import (
	"net/http"

	"github.com/fmaulll/lectureon/initializers"
	"github.com/fmaulll/lectureon/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx *gin.Context) {
	var body struct {
		Email    string
		Username string
		Password string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read body"})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to has password "})

		return
	}

	user := models.User{Email: body.Email, Username: body.Username, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User " + body.Username + " created!"})
}

func Login(ctx *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read body"})

		return
	}

}
