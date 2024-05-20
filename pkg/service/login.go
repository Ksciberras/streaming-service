package service

import (
	"fmt"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"video_stream/pkg/model"
)

type LoginService struct {
	db *gorm.DB
}

func NewLoginService(db *gorm.DB) *LoginService {
	return &LoginService{
		db: db,
	}
}

func (ls *LoginService) hashAndSalt(password string) string {
	log.Info("Hashing and salting password")

	passwordInBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordInBytes, bcrypt.MinCost)
	if err != nil {
		log.Error("Error while hashing and salting password")
		panic(err)
	}

	log.Info("Password hashed and salted")
	return string(hashedPassword)
}

func (ls *LoginService) SignUp(username string, password string) {
	log.Info("Signing up user")

	hashedPassword := ls.hashAndSalt(password)

	user := model.User{
		Username: username,
		Password: hashedPassword,
	}

	db := ls.db.Create(&user)
	if db.Error != nil {
		log.Error(fmt.Sprintf("Error while signing up user: %v", db.Error))
		panic(db.Error)
	}

	log.Info("User signed up")
}
