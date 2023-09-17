package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/janeczku/go-spinner"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Input struct {
	startTime   string
	endTime     string
	videoPath   string
	outputPath  string
	defaultPath string
}

func getInputs() *Input {
	defaultPath, err := os.Getwd()
	errorHandler(err, "Could not get the current directory")

	startTime := flag.String("start", "00:00:00", "Start time")
	endTime := flag.String("end", "", "End time")
	videoPath := flag.String("path", "", "URL / The path to a video file on your local machine.")
	outputPath := flag.String("out", defaultPath, "The path for the generated GIF.")

	flag.Parse()

	if len(*videoPath) < 1 {
		log.Fatal("Err: The path to video file is required")
	}

	if *endTime == "" {
		log.Fatal("Err: End time is required")
	}

	return &Input{
		startTime:   *startTime,
		endTime:     *endTime,
		videoPath:   *videoPath,
		outputPath:  *outputPath,
		defaultPath: defaultPath,
	}
}

func cropVideo(videoPath string, startTime string, endTime string) string {
	// Get directory to create the cropped video
	croppedDir, err := os.Getwd()
	errorHandler(err, "Could not read current directory")
	croppedVideoOutputPath := fmt.Sprintf("%s/cropped.mp4", croppedDir)

	// ffmpeg -ss 00:01:00 -to 00:02:00 -i input.mp4 -c copy output.mp4
	stream := ffmpeg.Input(videoPath, ffmpeg.KwArgs{"ss": startTime, "to": endTime})
	err = stream.Output(croppedVideoOutputPath, ffmpeg.KwArgs{"c": "copy", "v": "quiet"}).OverWriteOutput().ErrorToStdOut().Silent(true).Run()
	errorHandler(err, "Could not crop the video")

	return croppedVideoOutputPath
}

func videoToFrames(videoPath string) string {
	currentDir, err := os.Getwd()
	errorHandler(err, "Could not read current directory")

	framesDir := fmt.Sprintf("%s/frames", currentDir)
	if _, err := os.Stat(framesDir); os.IsNotExist(err) {
		err = os.Mkdir(framesDir, 0755)
		errorHandler(err, fmt.Sprintf("Could not create frames directory at %s", framesDir))
	}

	framesOutputImgPath := fmt.Sprintf("%s/img%%03d.png", framesDir)
	stream := ffmpeg.Input(videoPath)
	err = stream.Output(framesOutputImgPath, ffmpeg.KwArgs{"vf": "fps=15", "v": "quiet"}).OverWriteOutput().ErrorToStdOut().Silent(true).Run()
	errorHandler(err, "Could not convert cropped video to frames")

	return framesDir
}

func framesToGIF(framesPath string, outputPath string) string {
	framesImgPath := fmt.Sprintf("%s/img%%03d.png", framesPath)
	outputFilePath := fmt.Sprintf("%s/output.gif", outputPath)

	stream := ffmpeg.Input(framesImgPath, ffmpeg.KwArgs{"f": "image2", "framerate": "15", "loop": "0"})
	err := stream.Output(outputFilePath, ffmpeg.KwArgs{"v": "quiet"}).OverWriteOutput().ErrorToStdOut().Silent(true).Run()
	errorHandler(err, "Could not create GIF using frames")

	return outputFilePath
}

func cleanUp(dir string, filePath string) {
	err := os.RemoveAll(dir)
	errorHandler(err, fmt.Sprintf("Could not delete %s", dir))

	err = os.Remove(filePath)
	errorHandler(err, fmt.Sprintf("Could not delete %s", filePath))
}

func fetchVideo(url string, pathToSave string) string {
	videoExtension, err := getFileExtensionFromUrl(url)
	errString := "Error fetching resource"

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

	out, err := os.Create(fmt.Sprintf("%s/input.%s", pathToSave, videoExtension))
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

func main() {
	// Get inputs from the command line
	inputData := getInputs()
	videoPath := inputData.videoPath

	s := spinner.NewSpinner("Fetching video from URL...")

	// Check if video path is a URL
	if isUrl(inputData.videoPath) {
		s.Start()
		videoPath = fetchVideo(inputData.videoPath, inputData.defaultPath)
		s.Stop()
		fmt.Println("✓ Video fetched")
	}

	s = spinner.NewSpinner("Cooking up the GIF...")
	s.SetCharset([]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"})

	s.Start()
	// FFMPEG work
	croppedOutput := cropVideo(videoPath, inputData.startTime, inputData.endTime)
	framesOutput := videoToFrames(croppedOutput)
	gifPath := framesToGIF(framesOutput, inputData.outputPath)

	// Cleanup
	cleanUp(framesOutput, croppedOutput)

	s.Stop()
	fmt.Println("✓ Your GIF: ", gifPath)
}
