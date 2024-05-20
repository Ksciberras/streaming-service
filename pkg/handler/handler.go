package handler

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

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
	router.GET("/logout", h.HandleLogout)
}

func isUserAuthenticated(c *gin.Context, tokenKey string) bool {
	token, err := c.Cookie(tokenKey)
	if err != nil {
		log.Error(fmt.Sprintf("Error while getting token: %v", err))
		return false
	}
	err = service.VerifyToken(token)
	if err != nil {
		log.Error(fmt.Sprintf("Error while verifying token: %v", err))
		return false

	}
	return true
}

func (h *Handler) HandleTemplate(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Main website",
	})
}

func (h *Handler) HandleVideoTemplate(c *gin.Context) {
	isAuth := isUserAuthenticated(c, "token")
	if !isAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("uuid")
	video := h.videoService.FetchVideoById(id)

	c.HTML(200, "video.html", gin.H{
		"title": video.Title,
		"url":   h.serverUrl,
		"path":  fmt.Sprintf("%s/video/%s/%s.m3u8", h.serverUrl, id, id),
	})
}

func (h *Handler) AllVideos(c *gin.Context) {
	isAuth := isUserAuthenticated(c, "token")
	if !isAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Info("Fetching videos")
	videos := h.videoService.FetchVideos()
	c.JSON(200, videos)
}

func (h *Handler) ProcessDir(c *gin.Context) {
	h.videoService.ProcessVideosInDirectoryForStream("./videos")
}
