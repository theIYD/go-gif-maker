package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func errorHandler(err error, message string) {
	if err != nil && len(message) > 0 {
		log.Fatalf("Err: %s\n Trace: %s", message, err.Error())
	}
}

type Input struct {
	startTime  string
	endTime    string
	videoPath  string
	outputPath string
}

func getInputs() *Input {
	defaultPath, err := os.Getwd()
	errorHandler(err, "Could not get current directory")

	startTime := flag.String("start", "0", "Start time")
	endTime := flag.String("end", "0", "End time")
	videoPath := flag.String("path", defaultPath, "Video path")
	outputPath := flag.String("out", defaultPath, "Output path")

	flag.Parse()

	return &Input{
		startTime:  *startTime,
		endTime:    *endTime,
		videoPath:  *videoPath,
		outputPath: *outputPath,
	}
}

func cropVideo(videoPath string, startTime string, endTime string) string {
	// Get directory to create the cropped video
	croppedDir, err := os.Getwd()
	errorHandler(err, "Could not read current directory")

	inputVideoPath := fmt.Sprintf("%s/input.mp4", videoPath)
	croppedVideoOutputPath := fmt.Sprintf("%s/cropped.mp4", croppedDir)

	// ffmpeg -ss 00:01:00 -to 00:02:00 -i input.mp4 -c copy output.mp4
	stream := ffmpeg.Input(inputVideoPath, ffmpeg.KwArgs{"ss": startTime, "to": endTime})
	err = stream.Output(croppedVideoOutputPath, ffmpeg.KwArgs{"c": "copy"}).OverWriteOutput().ErrorToStdOut().Silent(true).Run()
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
	err = stream.Output(framesOutputImgPath, ffmpeg.KwArgs{"vf": "fps=15"}).OverWriteOutput().ErrorToStdOut().Silent(true).Run()
	errorHandler(err, "Could not convert cropped video to frames")

	return framesDir
}

func framesToGIF(framesPath string, outputPath string) string {
	framesImgPath := fmt.Sprintf("%s/img%%03d.png", framesPath)
	outputFilePath := fmt.Sprintf("%s/output.gif", outputPath)

	stream := ffmpeg.Input(framesImgPath, ffmpeg.KwArgs{"f": "image2", "framerate": "15", "loop": "0"})
	err := stream.Output(outputFilePath).OverWriteOutput().ErrorToStdOut().Silent(true).Run()
	errorHandler(err, "Could not create GIF using frames")

	return outputFilePath
}

func cleanUp(dir string, filePath string) {
	err := os.RemoveAll(dir)
	errorHandler(err, fmt.Sprintf("Could not delete %s", dir))

	err = os.Remove(filePath)
	errorHandler(err, fmt.Sprintf("Could not delete %s", filePath))
}

func main() {
	// Get inputs from the command line
	inputData := getInputs()

	// Handle the output directory
	outputDir := fmt.Sprintf("%s/output", inputData.outputPath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		errorHandler(err, fmt.Sprintf("Could not create output directory at %s", outputDir))
	}

	croppedOutput := cropVideo(inputData.videoPath, inputData.startTime, inputData.endTime)
	framesOutput := videoToFrames(croppedOutput)
	gifPath := framesToGIF(framesOutput, outputDir)

	// Cleanup
	cleanUp(framesOutput, croppedOutput)

	fmt.Println("Your GIF: ", gifPath)
}
