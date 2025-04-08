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
