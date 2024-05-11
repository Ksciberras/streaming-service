package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("index.html")
	router.GET("/", handleTemplate)
	router.StaticFS("/video", http.Dir("./hls"))

	router.Run(":8080")
}

func handleTemplate(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Main website",
	})
}

func newHlsVideo(inputPath string, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-profile:v", "baseline", "-level", "3.0", "-s", "640x360", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", outputPath)
	err := cmd.Run()
	return err
}

func streamVideo(c *gin.Context) {
	filePath := "Keith Vs Gil Catarino.mp4"
	outputPath := "./hls/keithgil/keithVGil.m3u8"
	err := newHlsVideo(filePath, outputPath)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusInternalServerError, "Error creating hls video")
	}
}
