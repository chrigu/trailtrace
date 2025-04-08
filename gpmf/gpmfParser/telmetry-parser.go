package gpmfParser

import (
	"fmt"
	"io"

	"gopro/mp4"
	"gopro/parser"
)

type TimeSample struct {
	TimeStamp int64
}

type TimedGPS struct {
	parser.GPS9
	TimeSample
}

// todo: rename
type TimedGyro struct {
	parser.Gyroscope
	TimeSample
}

type TimedFace struct {
	parser.Face
	TimeSample
}

func ExtractTelemetryData(file io.ReadSeeker, printTree bool) ([]TimedGPS, []TimedGyro, []TimedFace) {
	data, telemetryMetadata := mp4.ExtractTelemetryFromMp4(file)
	klvs := parser.ParseGPMF(data)

	if printTree {
		PrintTree(klvs, "")
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

func AddTimestampsToGPSData(gpsData []parser.GPS9, telemetryMetadata *mp4.TelemetryMetadata) []TimedGPS {
	var TimedGPSs []TimedGPS
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for i := 0; i < int(timeToSample.SampleCount); i++ {

			if sampleIndex >= uint32(len(gpsData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			TimedGPSs = append(TimedGPSs, TimedGPS{GPS9: gpsData[sampleIndex], TimeSample: TimeSample{TimeStamp: sampleTime}})
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return TimedGPSs
}

// refactor with AddTimestampsToGPSData
func AddTimestampsToFaceData(faceData [][]parser.Face, telemetryMetadata *mp4.TelemetryMetadata) []TimedFace {
	var TimedFaces []TimedFace
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for i := 0; i < int(timeToSample.SampleCount); i++ {
			if sampleIndex >= uint32(len(faceData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			for _, face := range faceData[sampleIndex] {
				TimedFaces = append(TimedFaces, TimedFace{Face: face, TimeSample: TimeSample{TimeStamp: sampleTime}})
			}
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return TimedFaces
}

// todo: refactor
func AddTimestampsToGyroDataWithDownsample(
	gyroData [][]parser.Gyroscope,
	telemetryMetadata *mp4.TelemetryMetadata,
	downsampleIntervalMs uint32,
) []TimedGyro {
	var TimedGyros []TimedGyro
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	var accumulatedGyro parser.Gyroscope
	var accumulatedTime int64 = 0
	var count uint32 = 0
	var lastSampleScaleTime int64 = 0

	// Precompute factor to avoid repeated calculation
	downsampleScaleThreshold := int64(telemetryMetadata.TimeScale * downsampleIntervalMs / 1000)

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range int(timeToSample.SampleCount) {
			if sampleIndex >= uint32(len(gyroData)) {
				break
			}

			currentTimedGyros := gyroData[sampleIndex]
			sampleCount := uint32(len(currentTimedGyros))

			for _, gyro := range currentTimedGyros {
				accumulatedTime += int64(sampleScaleTime)

				// Accumulate gyro values
				accumulatedGyro.X += gyro.X
				accumulatedGyro.Y += gyro.Y
				accumulatedGyro.Z += gyro.Z
				count++

				// Check if enough time has passed to downsample
				if int64(sampleScaleTime)-lastSampleScaleTime >= downsampleScaleThreshold {
					avgGyro := averageGyro(accumulatedGyro, count)
					avgTime := calculateAverageTime(telemetryMetadata.CreationTime, accumulatedTime, count, telemetryMetadata.TimeScale)

					TimedGyros = append(TimedGyros, TimedGyro{
						Gyroscope: avgGyro,
						TimeSample: TimeSample{
							TimeStamp: avgTime,
						},
					})

					// Reset accumulators
					accumulatedGyro = parser.Gyroscope{}
					lastSampleScaleTime = int64(sampleScaleTime)
					accumulatedTime = 0
					count = 0
				}

				// Increment scaled time based on sample delta
				sampleScaleTime += timeToSample.SampleDelta / sampleCount
			}
			sampleIndex++
		}
	}

	return TimedGyros
}

// Helper: Compute average Gyroscope reading
func averageGyro(accumulated parser.Gyroscope, count uint32) parser.Gyroscope {
	return parser.Gyroscope{
		X: accumulated.X / float32(count),
		Y: accumulated.Y / float32(count),
		Z: accumulated.Z / float32(count),
	}
}

// Helper: Compute average timestamp
func calculateAverageTime(creationTime int64, accumulatedTime int64, count uint32, timeScale uint32) int64 {
	averageScaleTime := accumulatedTime / int64(count)
	return creationTime + 1000*(averageScaleTime/int64(timeScale))
}
