package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Username string `json:"username"`
	Password string
	Role     string `json:"role"`
}
