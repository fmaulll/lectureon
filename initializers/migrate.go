package initializers

import "github.com/fmaulll/lectureon/models"

func Migrate() {
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Image{})
}
