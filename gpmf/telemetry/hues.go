package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedColor struct {
	parser.Color
	TimeSample
}

func AddTimestampsToHueData(colorData []parser.Color, telemetryMetadata *mp4.TelemetryMetadata) []TimedColor {
	var timedHues []TimedColor
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for i := 0; i < int(timeToSample.SampleCount); i++ {

			if sampleIndex >= uint32(len(colorData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			timedHues = append(timedHues, TimedColor{Color: colorData[sampleIndex], TimeSample: TimeSample{TimeStamp: sampleTime}})
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return timedHues
}
