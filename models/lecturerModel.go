package models

import (
	"time"
)

type Lecturer struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	PostId    int64     `json:"postId"`
}
