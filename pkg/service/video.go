package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"video_stream/pkg/model"
)

type VideoService struct {
	db *gorm.DB
}
type status string

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{db: db}
}

func (vs *VideoService) CheckAndCreateDirectory(path string) status {
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

func (vs *VideoService) RunHlsScript(inputPath string, outputPath string) error {
	log.Info("Running hls script")
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-profile:v", "baseline", "-level", "3.0", "-s", "640x360", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (vs *VideoService) ProcessVideosInDirectoryForStream(dirPath string) {
	var wg sync.WaitGroup

	var l sync.Mutex

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Checking for hls directory")
	status := vs.CheckAndCreateDirectory("./hls")
	if status == "exists" {
		log.Info("Hls directory exists")
	}

	for _, file := range files {
		log.Info(file.Name(), "isDir", file.IsDir())
		if file.IsDir() && strings.HasSuffix(file.Name(), ".mp4") {
			// continue
		}
		wg.Add(1)
		go vs.ProcessVideoForStream(file.Name(), dirPath, &wg, &l)

	}
}

func (vs *VideoService) AddVideo(request model.VideoRequest, uuid string) {
	log.Info("Adding video to database")

	video := model.Video{Id: uuid, Title: request.Title, Path: request.Path}

	vs.db.Create(&video)

	log.Info(fmt.Sprintf("Added video with id: %s title: %s", uuid, request.Title))
}

func (vs *VideoService) FetchVideos() []model.Video {
	log.Info("Fetching videos")
	var videos []model.Video
	vs.db.Find(&videos)
	return videos
}

func (vs *VideoService) FetchVideoById(id string) model.Video {
	var video model.Video
	vs.db.Where("id = ?", fmt.Sprintf("%s", id)).First(&video)
	return video
}

func (vs *VideoService) ProcessVideoForStream(fileName string, dirPath string, wg *sync.WaitGroup, l *sync.Mutex) {
	defer wg.Done()

	log.Info("Processing video", fileName)

	l.Lock()
	uuid := uuid.New().String()
	log.Info("New UUID ", uuid)

	outputDirPath := fmt.Sprintf("./hls/%s", uuid)
	status := vs.CheckAndCreateDirectory(outputDirPath)
	if status == "exists" {
		log.Info("Directory exists")
		l.Unlock()
		return
	}
	filePath := fmt.Sprintf("./%s/%s", dirPath, fileName)
	outputPath := fmt.Sprintf("%s/%s.m3u8", outputDirPath, uuid)
	err := vs.RunHlsScript(filePath, outputPath)
	if err != nil {
		log.Fatal(err)

		l.Unlock()
	}

	log.Info("Adding video to database")

	request := model.VideoRequest{Title: fileName, Path: fmt.Sprintf("video/%s/%s.m3u8", uuid, uuid)}

	vs.AddVideo(request, uuid)

	l.Unlock()
}
