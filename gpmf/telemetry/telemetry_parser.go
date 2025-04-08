package telemetry

import (
	"fmt"
	"io"

	"gopro/mp4"
	"gopro/parser"
)

func ExtractTelemetryData(file io.ReadSeeker, printTree bool) ([]TimedGPS, []TimedGyro, []TimedFace) {
	data, telemetryMetadata := mp4.ExtractTelemetryFromMp4(file)
	klvs := parser.ParseGPMF(data)

	if printTree {
		parser.PrintTree(klvs, "")
	}

	fmt.Println("KLVs", len(klvs))
	gpsData := parser.ParseGPS9Data(klvs)
	gyroData := parser.ParseGyroscopeData(klvs)
	accData := parser.ParseAccelerometerData(klvs)
	faceData := parser.ParseFaceData(klvs)
	fmt.Println("GPS9 data:", len(gpsData), "Gyro data:", len(gyroData), "Acc data:", len(accData), "Face data:", len(faceData))
	flattenedGpsData := make([]parser.GPS9, 0)
	for _, gpsSlice := range gpsData {
		flattenedGpsData = append(flattenedGpsData, gpsSlice...)
	}

	// todo: refactor
	gpsDataSamples := AddTimestampsToGPSData(flattenedGpsData, &telemetryMetadata)
	gyroDataSamples := AddTimestampsToGyroDataWithDownsample(accData, &telemetryMetadata, 250)
	faceDataSamples := AddTimestampsToFaceData(faceData, &telemetryMetadata)
	return gpsDataSamples, gyroDataSamples, faceDataSamples
}
