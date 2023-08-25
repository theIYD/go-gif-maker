package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// get path to read images from
	currentDir, err := os.Getwd()
	errorHandler(err)

	path := flag.String("path", currentDir, "The path to images")
	flag.Parse()

	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		// create output dir
		err = os.Mkdir("output", 0755)
		errorHandler(err)
	}

	// read path directory
	err = os.Chdir(*path)
	errorHandler(err)

	dir, err := os.ReadDir(".")
	errorHandler(err)

	// get the file name
	filePattern := strings.Split(dir[0].Name(), "-")
	inputFileName := fmt.Sprintf("%s-%%d.png", filePattern[0])
	outputFilePath := fmt.Sprintf("../output/%s-final.gif", filePattern[0])

	err = ffmpeg.Input(inputFileName, ffmpeg.KwArgs{"f": "image2", "framerate": "10", "loop": "0"}).Output(outputFilePath).OverWriteOutput().ErrorToStdOut().Run()
	errorHandler(err)
}
