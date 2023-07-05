package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fmaulll/lectureon/initializers"
	"github.com/fmaulll/lectureon/models"
	"github.com/gin-gonic/gin"
)

func NewPost(ctx *gin.Context) {
	title := ctx.PostForm("title")
	subTitle := ctx.PostForm("subTitle")
	description := ctx.PostForm("description")
	videoUrl := ctx.PostForm("videoUrl")

	// Source
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	t := time.Now().Unix()

	stringTime := strconv.Itoa(int(t))

	filename := stringTime + filepath.Base(file.Filename)
	dst := "./images/" + filename
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	post := models.Post{Title: title, SubTitle: subTitle, Description: description, Image: ctx.Request.Host + "/images/" + filename, VideoUrl: videoUrl}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post created!"})
}
