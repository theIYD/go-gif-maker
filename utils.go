package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

const YOUTUBE = "youtube"

func errorHandler(err error, message string) {
	if err != nil && len(message) > 0 {
		log.Fatalf("Err: %s\n Trace: %s", message, err.Error())
	}
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isYoutubeUrl(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	if !strings.Contains(u.Host, YOUTUBE) {
		return "", nil
	}

	return u.Query().Get("v"), nil
}

func getFileExtensionFromUrl(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return "", errors.New("couldn't find a period to indicate a file extension")
	}
	return u.Path[pos+1 : len(u.Path)], nil
}

func fetchYoutubeVideo(videoId string, pathToSave string) string {
	errString := "Error fetching resource."
	client := youtube.Client{}

	video, err := client.GetVideo(videoId)
	if err != nil {
		errorHandler(err, errString)
		return ""
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		errorHandler(err, errString)
		return ""
	}
	defer stream.Close()

	path := fmt.Sprintf("%s/input.mp4", pathToSave)
	file, err := os.Create(path)
	if err != nil {
		errorHandler(err, errString)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		errorHandler(err, errString)
		return ""
	}

	return path
}

func fetchVideo(url string, pathToSave string) string {
	videoExtension, err := getFileExtensionFromUrl(url)
	errString := "Error fetching resource."

	if err != nil {
		errorHandler(err, errString)
		return ""
	}

	resp, err := http.Get(url)
	if err != nil {
		errorHandler(err, errString)
		return ""
	}
	defer resp.Body.Close()

	path := fmt.Sprintf("%s/input.%s", pathToSave, videoExtension)

	out, err := os.Create(path)
	if err != nil {
		errorHandler(err, errString)
		return ""
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		errorHandler(err, errString)
	}

	return path
}
