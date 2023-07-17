package services

import (
	"discord-bot/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Service struct {
	maxDurationInSeconds int
	fileDirectory        string
}

func NewService(maxDurationInSeconds int) *Service {
	_ = os.Mkdir("data", os.ModePerm)
	return &Service{
		maxDurationInSeconds: maxDurationInSeconds,
		fileDirectory:        "data",
	}
}

func (svc *Service) SearchYTAndDownload(query string) (*models.Media, error) {
	timeout := make(chan bool, 1)
	result := make(chan searchAndDownloadResult, 1)
	go func() {
		time.Sleep(60 * time.Second)
		timeout <- true
	}()
	go func() {
		result <- svc.doSearchYTAndDownload(query)
	}()
	select {
	case <-timeout:
		return nil, errors.New("timed out")
	case result := <-result:
		return result.Media, result.Error
	}
}

func (svc *Service) doSearchYTAndDownload(query string) searchAndDownloadResult {
	start := time.Now()
	youtubeDownloader, err := exec.LookPath("yt-dlp")
	if err != nil {
		return searchAndDownloadResult{Error: errors.New("yt-dlp not found in path")}
	} else {
		args := []string{
			fmt.Sprintf("ytsearch10:%s", strings.ReplaceAll(query, "\"", "")),
			"--extract-audio",
			"--audio-format", "mp3",
			"--no-playlist",
			"--match-filter", fmt.Sprintf("duration < %d & !is_live", svc.maxDurationInSeconds),
			"--max-downloads", "1",
			"--output", fmt.Sprintf("%s/%d-%%(id)s.opus", svc.fileDirectory, start.Unix()),
			"--quiet",
			"--print-json",
			"--ignore-errors", // Ignores unavailable videos
			"--no-color",
			"--no-check-formats",
		}
		log.Printf("yt-dlp %s", strings.Join(args, " "))
		cmd := exec.Command(youtubeDownloader, args...)
		if data, err := cmd.Output(); err != nil && err.Error() != "exit status 101" {
			return searchAndDownloadResult{Error: fmt.Errorf("failed to search and download audio: %s\n%s", err.Error(), string(data))}
		} else {
			videoMetadata := models.VideoMetadata{}
			err = json.Unmarshal(data, &videoMetadata)
			if err != nil {
				fmt.Println(string(data))
				return searchAndDownloadResult{Error: fmt.Errorf("failed to unmarshal video metadata: %w", err)}
			}
			return searchAndDownloadResult{
				Media: models.NewMedia(
					videoMetadata.Title,
					videoMetadata.Filename+".mp3",
					videoMetadata.Uploader,
					fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoMetadata.ID),
					videoMetadata.Thumbnail,
					int(videoMetadata.Duration),
				),
			}
		}
	}
}

type searchAndDownloadResult struct {
	Media *models.Media
	Error error
}

