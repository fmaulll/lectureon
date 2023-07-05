package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string
	SubTitle    string `json:"subTitle"`
	Description string
	Image       string
	VideoUrl    string `json:"videoUrl"`
}
