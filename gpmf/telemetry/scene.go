package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedScene struct {
	Scenes []parser.Scene
	TimeSample
}

func AddTimestampsToSceneData(sceneData [][]parser.Scene, telemetryMetadata *mp4.TelemetryMetadata) []TimedScene {
	var timedScenes []TimedScene
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range timeToSample.SampleCount {

			if sampleIndex >= uint32(len(sceneData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			timedScenes = append(timedScenes, TimedScene{Scenes: sceneData[sampleIndex], TimeSample: TimeSample{TimeStamp: sampleTime}})
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return timedScenes
}
