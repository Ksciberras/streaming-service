package service

import (
	"bytes"
	"fmt"
	"os"

	"github.com/charmbracelet/log"

	"video_stream/pkg/model"
)

func (vs *VideoService) InitM3uPlaylist(videoDirectory string) {
	db := vs.db
	url := os.Getenv("SERVER_URL")
	path := "/playlist.m3u"
	if _, err := os.Stat(path); os.IsExist(err) {
		log.Info("Directory exists")
		log.Info(fmt.Sprintf("Removing file: %s", path))

		err := os.Remove(path)
		if err != nil {
			log.Error(fmt.Sprintf("Error while removing file: %v", err))
			return
		}
		fmt.Println("File removed")
	}

	files, err := os.ReadDir(videoDirectory)
	if err != nil {
		log.Error(fmt.Sprintf("Error while reading directory: %v", err))
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("#EXTM3U\n")
	fileNames := []string{}
	for _, file := range files {
		log.Info(fmt.Sprintf("File: %s", file.Name()))
		fileNames = append(fileNames, file.Name())
	}
	var videos []model.Video
	db.Where("Id IN ?", fileNames).Find(&videos)
	if db.Error != nil {
		log.Error(fmt.Sprintf("Error while fetching videos: %v", db.Error))
		return
	}

	for index, video := range videos {
		log.Info(fmt.Sprintf("Video: %s", video.Title))
		buffer.WriteString(fmt.Sprintf("#EXTINF:-%s,%s\n", index, video.Title))
		buffer.WriteString(fmt.Sprintf("%s/video/%s/%s.m3u8\n", url, video.Id, video.Id))

	}

	err = os.WriteFile("playlist.m3u", buffer.Bytes(), 0644)
	if err != nil {
		log.Error(fmt.Sprintf("Error while writing file: %v", err))
		return
	}
	log.Info("Playlist created")
}
