package service

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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

func (ls *LoginService) SignUp(username string, password string) error {
	log.Info(fmt.Sprintf("Logging in user: %s", username))

	hashedPassword := ls.hashAndSalt(password)

	user := model.User{
		Id:       uuid.New().String(),
		Username: username,
		Password: hashedPassword,
	}

	db := ls.db.Create(&user)
	if db.Error != nil {
		log.Error(fmt.Sprintf("Error while signing up user: %v", db.Error))
		return db.Error
	}

	log.Info("User signed up")
	return nil
}

type PasswordStatus string

const (
	valid   PasswordStatus = "valid"
	invalid PasswordStatus = "invalid"
	noValue PasswordStatus = "no-value"
)

func (ls *LoginService) CheckPassword(storedPassword string, loginRequestPassword string) PasswordStatus {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginRequestPassword))
	if err != nil {
		log.Error(fmt.Sprintf("Error while logging in user: %v", err))
		return "invalid"
	}
	return "valid"
}

func (ls *LoginService) Login(username string, password string) (PasswordStatus, model.User, error) {
	log.Info(fmt.Sprintf("Logging in user: %s", username))
	db := ls.db
	var user model.User
	db.First(&user, "username = ?", username)
	if db.Error != nil {
		log.Error(fmt.Sprintf("Error while logging in user: %v", db.Error))
		return "novalue", user, db.Error
	}

	passwordValidation := ls.CheckPassword(user.Password, password)

	return passwordValidation, user, nil
}
