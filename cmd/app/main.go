package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"video_stream/pkg/config"
	"video_stream/pkg/db"
	"video_stream/pkg/handler"
	"video_stream/pkg/service"
)

func main() {
	router := gin.Default()

	config.Cors(router)

	db := db.Connect()

	err := godotenv.Load("./pkg/config/.env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("SERVER_PORT")
	url := os.Getenv("SERVER_URL")

	videoService := service.NewVideoService(db)
	loginService := service.NewLoginService(db)
	handler := handler.NewHandler(videoService, url, loginService)
	handler.SetupRoutes(router)

	router.Run(port)
}
