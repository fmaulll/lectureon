package controllers

import (
	"fmt"
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
	// videoUrl := ctx.PostForm("videoUrl")

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	id, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	post := models.Post{Title: title, SubTitle: subTitle, Description: description, AuthorID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})

		return
	}

	for i := 1; ; i++ {
		file, err := ctx.FormFile("image" + strconv.Itoa(i))
		if err != nil {
			// No more files with the current key, break the loop
			break
		}

		t := time.Now().Unix()
		stringTime := strconv.Itoa(int(t))

		filename := stringTime + filepath.Base(file.Filename)

		dst := "./images/" + filename

		// Handle the uploaded file as needed
		// For example, save it to disk or process it in some way
		err = ctx.SaveUploadedFile(file, dst)
		if err != nil {
			// Handle the error
			ctx.String(http.StatusInternalServerError, "Failed to upload file")
			return
		}

		if err := ctx.SaveUploadedFile(file, dst); err != nil {
			ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		image := models.Image{CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: file.Filename, Url: scheme + "://" + ctx.Request.Host + "/images/" + filename, PostId: post.ID}

		if err := initializers.DB.Create(&image); err.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload image"})

			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "post created!"})
}

func GetAllPost(ctx *gin.Context) {

	type Post struct {
		ID          int64  `json:"id"`
		AuthorID    int64  `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Title       string `json:"title"`
		SubTitle    string `json:"subTitle"`
		Description string `json:"description"`
		Images      string `json:"images"`
		// VideoUrl    string `json:"videoUrl"`
	}

	var results []Post

	if err := initializers.DB.Table("users").Select("posts.id, users.id as author_id, users.first_name || ' ' || users.last_name as author_name, posts.title, posts.sub_title, posts.description, json_agg(images.url) as images").Joins("JOIN posts ON users.id = posts.author_id").Joins("LEFT JOIN images ON posts.id = images.post_id").Group("users.id, posts.id, images.post_id").Order("posts.created_at desc").Limit(10).Scan(&results); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "post not found!"})

		return
	}

	ctx.JSON(http.StatusOK, results)
}

func GetAllPostByAuthorId(ctx *gin.Context) {
	id := ctx.Param("id")

	var posts []models.Post

	if err := initializers.DB.Where("author_id = ?", id).Find(&posts); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "post not found!"})

		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func EditPost(ctx *gin.Context) {
	var body struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		SubTitle    string `json:"subTitle"`
		Description string `json:"description"`
		Image       string `json:"image"`
		VideoUrl    string `json:"videoUrl"`
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read body"})

		return
	}

	result := initializers.DB.Model(&models.Post{}).Find(&models.Post{}, body.ID).Select("title", "sub_title", "description", "image", "video_url", "updated_at").Updates(map[string]interface{}{"title": body.Title, "sub_title": body.SubTitle, "description": body.Description, "image": body.Image, "video_url": body.VideoUrl, "updated_at": time.Now()})

	if result.Error != nil {
		fmt.Println(result.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update post!"})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "successfully update post!"})
}
