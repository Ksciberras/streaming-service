package handler

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"video_stream/pkg/model"
	"video_stream/pkg/service"
)

func (h *Handler) handleLogin(c *gin.Context) {
	var signIn model.LoginRequest
	err := c.BindJSON(&signIn)
	if err != nil {
		log.Error(fmt.Sprintf("Error while binding json: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Info(fmt.Sprintf("Logging in user: %s %s ", signIn.Username, signIn.Password))

	passwordStatus, user, loginError := h.loginService.Login(signIn.Username, signIn.Password)
	if loginError != nil {
		log.Error(fmt.Sprintf("Error while logging in user: %v", loginError))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while logging in user"})
		return
	}
	if passwordStatus == "invalid" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if passwordStatus == "valid" {
		token, err := service.CreateToken(user.Username)
		if err != nil {
			log.Error(fmt.Sprintf("Error while creating token: %v", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating token"})
			return
		}
		c.SetCookie("token", token, 60*60*24, "/", "localhost", false, true)
		c.JSON(200, gin.H{"message": fmt.Sprintf("User: %s succsefully logged in", signIn.Username), "user": user, "token": token})
		return
	}
}

func (h *Handler) HandleLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "User logged out"})
}

func (h *Handler) handleSignUp(c *gin.Context) {
	var signUp model.LoginRequest

	err := c.BindJSON(&signUp)
	if err != nil {
		log.Error(fmt.Sprintf("Error while binding json: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Info(fmt.Sprintf("Signing up user: %s %s ", signUp.Username, signUp.Password))

	loginErr := h.loginService.SignUp(signUp.Username, signUp.Password)
	if loginErr != nil {
		log.Error(fmt.Sprintf("Error while signing up user: %v", loginErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while signing up user"})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("User: %s succsefully signed up", signUp.Username)})
}
