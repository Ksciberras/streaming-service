package handler

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"video_stream/pkg/model"
	"video_stream/pkg/service"
)

type Handler struct {
	loginService *service.LoginService
	videoService *service.VideoService
	serverUrl    string
}

func NewHandler(videoService *service.VideoService, serverUrl string, loginService *service.LoginService) *Handler {
	return &Handler{videoService: videoService, serverUrl: serverUrl, loginService: loginService}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.LoadHTMLFiles("../../web/index.html", "../../web/video.html")
	router.StaticFS("/video", http.Dir("./hls"))

	router.LoadHTMLFiles("./web/index.html", "./web/video.html")
	router.StaticFS("../../videos", http.Dir("./hls"))
	router.GET("/", h.HandleTemplate)
	router.GET("/v/:uuid", h.HandleVideoTemplate)
	router.GET("/videos", h.AllVideos)
	router.GET("/dir", h.ProcessDir)

	router.POST("/signup", h.handleSignUp)
	router.POST("/login", h.handleLogin)
}

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
		c.JSON(200, gin.H{"message": fmt.Sprintf("User: %s succsefully logged in", signIn.Username), "user": user})
		return
	}
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

func (h *Handler) HandleTemplate(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Main website",
	})
}

func (h *Handler) HandleVideoTemplate(c *gin.Context) {
	id := c.Param("uuid")
	video := h.videoService.FetchVideoById(id)

	c.HTML(200, "video.html", gin.H{
		"title": video.Title,
		"url":   h.serverUrl,
		"path":  fmt.Sprintf("%s/video/%s/%s.m3u8", h.serverUrl, id, id),
	})
}

func (h *Handler) AllVideos(c *gin.Context) {
	log.Info("Fetching videos")
	videos := h.videoService.FetchVideos()
	c.JSON(200, videos)
}

func (h *Handler) ProcessDir(c *gin.Context) {
	h.videoService.ProcessVideosInDirectoryForStream("./videos")
}
