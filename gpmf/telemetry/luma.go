package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedLuma struct {
	parser.Luma
	TimeSample
}

func AddTimestampsToLumaData(lumaData []parser.Luma, telemetryMetadata *mp4.TelemetryMetadata) []TimedLuma {
	var timedLumas []TimedLuma
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range timeToSample.SampleCount {

			if sampleIndex >= uint32(len(lumaData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			timedLumas = append(timedLumas, TimedLuma{Luma: lumaData[sampleIndex], TimeSample: TimeSample{TimeStamp: sampleTime}})
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return timedLumas
}
