package main

import (
	"fmt"
	"os"

	"gopro/telemetry"
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
	gpsData, gyroData, faceData, _, _ := telemetry.ExtractTelemetryData(file, false)

	fmt.Printf("GPS Data: %v\n", gpsData)
	fmt.Printf("Gyro Data: %v\n", gyroData)
	fmt.Printf("Face Data: %v\n", faceData)
	// fmt.Printf("Luma Data: %v\n", lumaData)
	// fmt.Printf("Hues Data: %v\n", huesData)

	// Test()
}
