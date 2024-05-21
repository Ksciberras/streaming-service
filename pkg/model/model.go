package model

import (
	"time"
)

type Video struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title     string    `json:"title"`
	Path      string    `json:"path"`
}

type VideoRequest struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}
type User struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
