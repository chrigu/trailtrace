package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedGPS struct {
	parser.GPS9
	TimeSample
}

func AddTimestampsToGPSData(gpsData []parser.GPS9, telemetryMetadata *mp4.TelemetryMetadata) []TimedGPS {
	var timedGPSs []TimedGPS
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range timeToSample.SampleCount {
			if sampleIndex >= uint32(len(gpsData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			timedGPSs = append(timedGPSs, TimedGPS{GPS9: gpsData[sampleIndex], TimeSample: TimeSample{TimeStamp: sampleTime}})
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return timedGPSs
}
