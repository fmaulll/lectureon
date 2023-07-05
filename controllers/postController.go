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

	post := models.Post{Title: title, SubTitle: subTitle, Description: description, Image: ctx.Request.Host + "/images/" + filename, VideoUrl: videoUrl, AuthorID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post created!"})
}

func GetAllPost(ctx *gin.Context) {

	// type dataResult struct {
	// 	ID          int64  `json:"id"`
	// 	Title       string `json:"title"`
	// 	SubTitle    string `json:"subTitle"`
	// 	Description string `json:"description"`
	// 	Image       string `json:"image"`
	// 	VideoUrl    string `json:"videoUrl"`
	// 	AuthorID    int64  `json:"authorId"`
	// 	AuthorName  string `json:"authorName"`
	// }

	// var postResult []dataResult

	var posts []models.Post
	result := initializers.DB.Find(&posts)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "post not found!"})

		return
	}

	// for i := range posts {
	// 	var author *models.User

	// 	if result := initializers.DB.First(&author, "id = ?", posts[i].AuthorID); result.Error != nil {
	// 		ctx.JSON(http.StatusNotFound, gin.H{"message": "post not found!"})

	// 		return
	// 	}

	// 	postResult[i].ID = posts[i].ID
	// 	postResult[i].AuthorID = posts[i].AuthorID
	// 	postResult[i].AuthorName = author.FirstName + author.LastName
	// 	postResult[i].Title = posts[i].Title
	// 	postResult[i].SubTitle = posts[i].SubTitle
	// 	postResult[i].Description = posts[i].Description
	// 	postResult[i].Image = posts[i].Image
	// 	postResult[i].VideoUrl = posts[i].VideoUrl
	// }

	ctx.JSON(http.StatusOK, gin.H{"results": posts})
	// ctx.JSON(http.StatusOK, gin.H{"results": postResult})
}
