package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedHue struct {
	Hues []parser.Hue
	TimeSample
}

func AddTimestampsToHueData(hueData [][]parser.Hue, telemetryMetadata *mp4.TelemetryMetadata) []TimedHue {
	var timedHues []TimedHue
	var sampleScaleTime uint32 = 0

	for _, hues := range hueData {
		sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
		timedHues = append(timedHues, TimedHue{
			Hues:       hues,
			TimeSample: TimeSample{TimeStamp: sampleTime},
		})
		sampleScaleTime += telemetryMetadata.TimeToSamples[0].SampleDelta
	}

	return timedHues
}
