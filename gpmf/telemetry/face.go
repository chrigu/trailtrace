package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedFace struct {
	parser.Face
	TimeSample
}

// refactor with AddTimestampsToGPSData
func AddTimestampsToFaceData(faceData [][]parser.Face, telemetryMetadata *mp4.TelemetryMetadata) []TimedFace {
	var TimedFaces []TimedFace
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range timeToSample.SampleCount {
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
