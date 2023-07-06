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
	authorId := ctx.PostForm("id")
	title := ctx.PostForm("title")
	subTitle := ctx.PostForm("subTitle")
	description := ctx.PostForm("description")
	videoUrl := ctx.PostForm("videoUrl")

	id, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

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
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	post := models.Post{Title: title, SubTitle: subTitle, Description: description, Image: scheme + "://" + ctx.Request.Host + "/images/" + filename, VideoUrl: videoUrl, AuthorID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post created!"})
}

func GetAllPost(ctx *gin.Context) {

	var result []struct {
		ID          int64  `json:"id"`
		AuthorID    int64  `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Title       string `json:"title"`
		SubTitle    string `json:"subTitle"`
		Description string `json:"description"`
		Image       string `json:"image"`
		VideoUrl    string `json:"videoUrl"`
	}

	if err := initializers.DB.Table("users").Select("posts.id, users.id as author_id, users.first_name || ' ' || users.last_name as author_name, posts.title, posts.sub_title, posts.description, posts.image, posts.video_url").Joins("JOIN posts ON users.id = posts.author_id").Scan(&result); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "post not found!"})

		return
	}

	ctx.JSON(http.StatusOK, result)
}
