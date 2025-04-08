package gpmfParser

import (
	"fmt"
	"io"

	"gopro/mp4"
)

type TimeSample struct {
	TimeStamp int64
}

type TimedGPS struct {
	GPS9
	TimeSample
}

// todo: rename
type TimedGyro struct {
	Gyroscope
	TimeSample
}

type TimedFace struct {
	Face
	TimeSample
}

func ExtractTelemetryData(file io.ReadSeeker, printTree bool) ([]TimedGPS, []TimedGyro, []TimedFace) {
	data, telemetryMetadata := mp4.ExtractTelemetryFromMp4(file)
	klvs := ParseGPMF(data)

	if printTree {
		PrintTree(klvs, "")
	}

	fmt.Println("KLVs", len(klvs))
	gpsData := parseGPS9Data(klvs)
	gyroData := parseGyroscopeData(klvs)
	accData := parseAccelerometerData(klvs)
	faceData := parsecFaceData(klvs)
	fmt.Println("GPS9 data:", len(gpsData), "Gyro data:", len(gyroData), "Acc data:", len(accData), "Face data:", len(faceData))
	flattenedGpsData := make([]GPS9, 0)
	for _, gpsSlice := range gpsData {
		flattenedGpsData = append(flattenedGpsData, gpsSlice...)
	}

	// todo: refactor
	gpsDataSamples := assignTimestampsToGps(flattenedGpsData, &telemetryMetadata)
	gyroDataSamples := assignTimestampsToGyroWithAverage(accData, &telemetryMetadata, 250)
	faceDataSamples := assignTimestampsToFace(faceData, &telemetryMetadata)
	return gpsDataSamples, gyroDataSamples, faceDataSamples
}

func assignTimestampsToGps(gpsData []GPS9, telemetryMetadata *mp4.TelemetryMetadata) []TimedGPS {
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

// refactor with assignTimestampsToGps
func assignTimestampsToFace(faceData [][]Face, telemetryMetadata *mp4.TelemetryMetadata) []TimedFace {
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
func assignTimestampsToGyroWithAverage(
	gyroData [][]Gyroscope,
	telemetryMetadata *mp4.TelemetryMetadata,
	downsampleIntervalMs uint32,
) []TimedGyro {
	var TimedGyros []TimedGyro
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	var accumulatedGyro Gyroscope
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
					accumulatedGyro = Gyroscope{}
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
func averageGyro(accumulated Gyroscope, count uint32) Gyroscope {
	return Gyroscope{
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
