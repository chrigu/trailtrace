package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedGPS struct {
	parser.GPS9
	TimeSample
}

func AddTimestampsToGPSData(gpsData [][]parser.GPS9, telemetryMetadata *mp4.TelemetryMetadata) []TimedGPS {
	var timedGPSs []TimedGPS
	var sampleIndex uint32 = 0
	var baseSampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range timeToSample.SampleCount {
			if sampleIndex >= uint32(len(gpsData)) {
				break
			}

			for i, gps := range gpsData[sampleIndex] {
				// Calculate time delta relative to the start of this time-to-sample entry
				scaledTimeDelta := float64(i) / float64(len(gpsData[sampleIndex])) * float64(timeToSample.SampleDelta)
				timeDelta := float64(baseSampleScaleTime+uint32(scaledTimeDelta)) * 1000.0 / float64(telemetryMetadata.TimeScale)
				sampleTime := telemetryMetadata.CreationTime + int64(timeDelta)
				timedGPSs = append(timedGPSs, TimedGPS{GPS9: gps, TimeSample: TimeSample{TimeStamp: sampleTime}})
			}
			sampleIndex++
			baseSampleScaleTime += timeToSample.SampleDelta
		}
	}

	return timedGPSs
}
