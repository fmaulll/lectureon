package models

import (
	"time"
)

type Post struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Title       string    `json:"title"`
	SubTitle    string    `json:"subTitle"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	VideoUrl    string    `json:"videoUrl"`
	AuthorID    int64     `json:"authorId"`
}
