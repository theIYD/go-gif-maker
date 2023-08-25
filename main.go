package main

import (
	"flag"
	"os"
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

	// read path directory
	dir, err := os.ReadDir(*path)
	errorHandler(err)

	var fileNames []string

	for _, file := range dir {
		fileNames = append(fileNames, file.Name())
	}
}
