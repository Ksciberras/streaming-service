package model

import (
	"time"
)

type Video struct {
	Id        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Path      string
}

type VideoRequest struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}
