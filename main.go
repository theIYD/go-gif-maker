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

/*
	- take in a dirpath (video)/dirpath (images)/video url for creating a gif
	- in case of video, take in start, end time for crop
	- use ffmpeg to get each image sequence and store
	- use below logic to read the images and generate a gif
*/

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
	err = stream.Output(croppedVideoOutputPath, ffmpeg.KwArgs{"c": "copy"}).OverWriteOutput().ErrorToStdOut().Run()
	errorHandler(err, "Could not crop the video")

	return croppedVideoOutputPath
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
	fmt.Println("Cropped video output", croppedOutput)

	// // read path directory
	// err = os.Chdir(*path)
	// errorHandler(err)

	// dir, err := os.ReadDir(".")
	// errorHandler(err)

	// // get the file name
	// filePattern := strings.Split(dir[0].Name(), "-")
	// inputFileName := fmt.Sprintf("%s-%%d.png", filePattern[0])
	// outputFilePath := fmt.Sprintf("../output/%s-final.gif", filePattern[0])

	// err = ffmpeg.Input(inputFileName, ffmpeg.KwArgs{"f": "image2", "framerate": "10", "loop": "0"}).Output(outputFilePath).OverWriteOutput().ErrorToStdOut().Run()
	// errorHandler(err)
}
