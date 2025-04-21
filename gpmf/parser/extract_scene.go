package parser

import (
	"gopro/internal"
)

type FourCCScene string

const (
	SceneSNOW       FourCCScene = "SNOW"
	SceneURBAN      FourCCScene = "URBA"
	SceneINDOOR     FourCCScene = "INDO"
	SceneWATER      FourCCScene = "WATR"
	SceneVEGETATION FourCCScene = "VEGE"
	SceneBEACH      FourCCScene = "BEAC"
)

type Scene struct {
	Type FourCCScene
	Prob float32
}

func ParseSceneData(klvs []KLV) [][]Scene {
	return extractSensorData(klvs,
		"Scene classification[[CLASSIFIER_FOUR_CC,prob], ...]",
		extractSceneData)
}

func extractSceneData(klv KLV) []Scene {

	// todo: extract types dynamically
	var format string = ""
	var payload []byte = make([]byte, 0)
	var repeat uint32 = 0

	for _, child := range klv.Children {
		// log("Processing child:", child.FourCC)

		switch child.FourCC {
		case "SCEN":
			internal.Log("SCEN found")
			payload = child.Payload
			repeat = child.Repeat
		case "TYPE":
			internal.Log("TYPE found")
			format = readPayload(child).(string)
		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	sceneRawData, err := parseDynamicStructure(payload, format, repeat) // make easier, check type and make struct
	if err != nil {
		internal.Log("Error parsing dynamic structure:", err)
		return []Scene{} // Return empty slice on error
	}

	sceneValues := make([]Scene, len(sceneRawData))

	// Process each GPS value
	for i, rawData := range sceneRawData {
		if len(rawData) < 2 {
			internal.Log("Incomplete scene data at index %d", i)
			continue
		}

		// Extract and convert the values with proper type assertions
		var sceneType FourCCScene
		if str, ok := rawData[0].(string); ok {
			sceneType = FourCCScene(str)
		} else if fourcc, ok := rawData[0].(FourCCScene); ok {
			sceneType = fourcc
		} else {
			internal.Log("Invalid type for scene data at index %d: %T", i, rawData[0])
			continue
		}
		sceneProb, ok2 := rawData[1].(float32)

		if !ok2 {
			internal.Log("Type assertion failed for scene data at index %d", i)
			continue
		}

		sceneValues[i] = Scene{
			Type: sceneType,
			Prob: sceneProb,
		}
	}

	internal.Log("Extracted %d scene values", len(sceneValues))
	return sceneValues
}
