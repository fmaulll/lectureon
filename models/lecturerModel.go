package models

import (
	"time"
)

type Lecturer struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	StudentID int64     `json:"studentId"`
}
