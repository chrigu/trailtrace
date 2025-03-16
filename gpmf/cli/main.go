package main

import (
	"fmt"
	"gopro/gpmfParser"
	"os"
)

func main() {
	// Open MP4 file for reading

	if len(os.Args) < 2 {
		fmt.Println("Error: No file specified")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Extract metadata track from the MP4 file
	gpmfParser.ExtractTelemetryDataFromMp4(file)
	// Test()
}
