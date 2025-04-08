package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedGyro struct {
	parser.Gyroscope
	TimeSample
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
