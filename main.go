package main

import (
	"fmt"
	"os"

)


func main() {
	// Open MP4 file for reading
	

	file, err := os.Open("../GX010025.mp4")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Extract metadata track from the MP4 file
	ExtractTelemetryData(file)
  // Test()
}

