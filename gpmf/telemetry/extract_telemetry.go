package telemetry

import (
	"fmt"
	"io"

	"gopro/mp4"
	"gopro/parser"
)

func ExtractTelemetryData(file io.ReadSeeker, printTree bool) ([]TimedGPS, []TimedGyro, []TimedFace, []TimedLuma, []TimedHue, []TimedScene) {
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
	lumaData := parser.ParseLumaData(klvs)
	hueData := parser.ParseHueData(klvs)
	sceneData := parser.ParseSceneData(klvs)
	fmt.Println("GPS9 data:", len(gpsData), "Gyro data:", len(gyroData), "Acc data:", len(accData), "Face data:", len(faceData), "luma data:", len(lumaData), "hue data:", len(hueData), "scene data:", len(sceneData))

	// refactor
	flattenedGpsData := make([]parser.GPS9, 0)
	for _, gpsSlice := range gpsData {
		flattenedGpsData = append(flattenedGpsData, gpsSlice...)
	}

	flattenedLumaData := make([]parser.Luma, 0)
	for _, lumaSlice := range lumaData {
		flattenedLumaData = append(flattenedLumaData, lumaSlice...)
	}

	flattenedSceneData := make([]parser.Scene, 0)
	for _, sceneSlice := range sceneData {
		flattenedSceneData = append(flattenedSceneData, sceneSlice...)
	}

	// todo: refactor
	timedGpsData := AddTimestampsToGPSData(flattenedGpsData, &telemetryMetadata)
	timedGyroData := AddTimestampsToGyroDataWithDownsample(accData, &telemetryMetadata, 250)
	timedFaceData := AddTimestampsToFaceData(faceData, &telemetryMetadata)
	timedLumaData := AddTimestampsToLumaData(flattenedLumaData, &telemetryMetadata)
	timedHueData := AddTimestampsToHueData(hueData, &telemetryMetadata)
	timedSceneData := AddTimestampsToSceneData(flattenedSceneData, &telemetryMetadata)
	return timedGpsData, timedGyroData, timedFaceData, timedLumaData, timedHueData, timedSceneData
}
