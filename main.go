package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func errorHandler(err error, message string) {
	if err != nil && len(message) > 0 {
		log.Fatalf("Err: %s\n Trace: %s", message, err.Error())
	}
}

type Input struct {
	startTime   string
	endTime     string
	videoPath   string
	outputPath  string
	defaultPath string
}

func getInputs() *Input {
	defaultPath, err := os.Getwd()
	errorHandler(err, "Could not get current directory")

	startTime := flag.String("start", "0", "Start time")
	endTime := flag.String("end", "0", "End time")
	videoPath := flag.String("path", defaultPath, "URL / The path to a video file on your local machine.")
	outputPath := flag.String("out", defaultPath, "Define the path for the GIF.")

	flag.Parse()

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
	videoPath := fmt.Sprintf("%s/input.mp4", inputData.videoPath)

	// Check if video path is a URL
	if isUrl(inputData.videoPath) {
		videoPath = fetchVideo(inputData.videoPath, inputData.defaultPath)
	}

	// Handle the output directory
	outputDir := fmt.Sprintf("%s/output", inputData.outputPath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		errorHandler(err, fmt.Sprintf("Could not create output directory at %s", outputDir))
	}

	croppedOutput := cropVideo(videoPath, inputData.startTime, inputData.endTime)
	framesOutput := videoToFrames(croppedOutput)
	gifPath := framesToGIF(framesOutput, outputDir)

	// Cleanup
	cleanUp(framesOutput, croppedOutput)

	fmt.Println("Your GIF: ", gifPath)
}
