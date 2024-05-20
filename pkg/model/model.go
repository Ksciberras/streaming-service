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
type User struct {
	Id        string `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
