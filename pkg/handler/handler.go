package handler

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"video_stream/pkg/service"
)

type Handler struct {
	videoService *service.VideoService
	serverUrl    string
}

func NewHandler(videoService *service.VideoService, serverUrl string) *Handler {
	return &Handler{videoService: videoService, serverUrl: serverUrl}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.LoadHTMLFiles("../../web/index.html", "../../web/video.html")
	router.StaticFS("/video", http.Dir("./hls"))

	router.GET("/", h.HandleTemplate)
	router.GET("/v/:uuid", h.HandleVideoTemplate)
	router.GET("/videos", h.AllVideos)
	router.GET("/dir", h.ProcessDir)
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
	service.ProcessVideosInDirectoryForStream("./videos")
}
