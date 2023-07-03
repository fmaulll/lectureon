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
		FirstName string
		LastName  string
		Email     string
		Username  string
		Role      string
		Password  string
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

	user := models.User{FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Username: body.Username, Role: body.Role, Password: string(hash)}

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

	var user *models.User
	initializers.DB.First(&user, "username = ?", body.Username)

	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})

		return
	}

	tokens, err := GenerateToken(user.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successfully",
		"access_token":  tokens["access_token"],
		"refresh_token": tokens["refresh_token"],
	})

}
