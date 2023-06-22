package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/fmaulll/lectureon/initializers"
	"github.com/fmaulll/lectureon/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create token"})

		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create token"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successfully",
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	})

}
