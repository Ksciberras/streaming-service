package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust the port according to your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.LoadHTMLFiles("index.html", "video.html")
	router.GET("/", handleTemplate)
	router.GET("/v/:uuid", handleVideoTemplate)
	router.GET("/videos", allVideos)
	router.GET("/dir", processDir)
	router.StaticFS("/video", http.Dir("./hls"))
	router.Run(":8080")
}

var url = "https://6341-77-71-159-208.ngrok-free.app"

func handleTemplate(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Main website",
	})
}

func handleVideoTemplate(c *gin.Context) {
	id := c.Param("uuid")
	db := connectDatabase()
	var video Video
	db.Where("id = ?", fmt.Sprintf("%s", id)).First(&video)

	c.HTML(200, "video.html", gin.H{
		"title": video.Title,
		"url":   url,
		"path":  fmt.Sprintf("%s/video/%s/%s.m3u8", url, id, id),
	})
}

func runHlsScript(inputPath string, outputPath string) error {
	log.Info("Running hls script")
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-profile:v", "baseline", "-level", "3.0", "-s", "640x360", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", outputPath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func processDir(c *gin.Context) {
	processVideosInDirectoryForStream("./videos")
}

type status string

func checkAndCreateDirectory(path string) status {
	const (
		exists  status = "exists"
		created status = "created"
	)
	log.Info("Checking and creating directory")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		log.Info("Directory created")
		return status(created)
	}
	log.Info("Directory already exists")
	return status(exists)
}

var l sync.Mutex

func processVideosInDirectoryForStream(dirPath string) {
	db := connectDatabase()
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Checking for hls directory")
	status := checkAndCreateDirectory("./hls")
	if status == "exists" {
		log.Info("Hls directory exists")
	}

	for _, file := range files {
		log.Info(file.Name(), "isDir", file.IsDir())
		if file.IsDir() && strings.HasSuffix(file.Name(), ".mp4") {
			continue
		}
		l.Lock()
		go processVideoForStream(file.Name(), db, dirPath)

	}
}

func processVideoForStream(fileName string, db *gorm.DB, dirPath string) {
	log.Info("Processing video", fileName)
	processVideosInDirectoryForStream("./hls")

	uuid := uuid.NewString()

	outputDirPath := fmt.Sprintf("./hls/%s", uuid)
	status := checkAndCreateDirectory(outputDirPath)
	if status == "exists" {
		log.Info("Directory exists")
		return
	}

	filePath := fmt.Sprintf("./%s/%s", dirPath, fileName)
	outputPath := fmt.Sprintf("%s/%s.m3u8", outputDirPath, uuid)
	err := runHlsScript(filePath, outputPath)
	if err != nil {
		log.Fatal(err)

		panic(err)
	}

	addVideo(db, VideoRequest{Title: fileName, Path: outputPath}, uuid)
	l.Unlock()
}

func addVideo(db *gorm.DB, request VideoRequest, uuid string) {
	log.Info("Adding video to database")

	video := Video{Id: uuid, Title: request.Title, Path: request.Path}

	db.Create(&video)

	log.Info(fmt.Sprintf("Added video with id: %s title: %s", uuid, request.Title))
}

func allVideos(c *gin.Context) {
	db := connectDatabase()
	videos := fetchVideos(db)
	c.JSON(200, videos)
}

func fetchVideos(db *gorm.DB) []Video {
	log.Info("Fetching videos")
	var videos []Video
	db.Find(&videos)
	return videos
}

func connectDatabase() *gorm.DB {
	log.Info("Connecting to database")

	db, err := gorm.Open(sqlite.Open("videos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to database")

	db.AutoMigrate(&Video{})

	return db
}
