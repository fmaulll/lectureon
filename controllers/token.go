package controllers

import (
	"os"
	"time"

	"github.com/fmaulll/lectureon/initializers"
	"github.com/fmaulll/lectureon/models"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(id uint) (map[string]string, error) {

	var user *models.User
	initializers.DB.First(&user, "id = ?", id)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}, nil
}
