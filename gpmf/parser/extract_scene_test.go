package parser

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"
)

func TestExtractSceneData(t *testing.T) {
	// 1. Prepare SCEN child (scene data)
	sceneValues := [][]interface{}{
		{FourCCScene("SNOW"), 0.8},
		{FourCCScene("URBA"), 0.6},
		{FourCCScene("INDO"), 0.4},
	}

	// Create a payload from the scene values
	var scenPayload []byte
	// Each entry needs 8 bytes: 4 for FourCC + 4 for float32
	for _, packet := range sceneValues {
		// Add FourCC (4 bytes)
		fourCC := packet[0].(FourCCScene)
		scenPayload = append(scenPayload, []byte(fourCC)...)

		// Add probability as float32 (4 bytes)
		probBytes := make([]byte, 4)
		prob := float32(packet[1].(float64))
		binary.BigEndian.PutUint32(probBytes, math.Float32bits(prob))
		scenPayload = append(scenPayload, probBytes...)
	}

	scenChild := KLV{
		FourCC:   "SCEN",
		DataType: int('B'), // Simplified for testing
		Repeat:   uint32(len(sceneValues)),
		DataSize: 1, // Simplified for testing
		Payload:  scenPayload,
	}

	// 2. Prepare TYPE child
	typeChild := KLV{
		FourCC:   "TYPE",
		DataType: int('c'), // char
		Repeat:   1,
		DataSize: 1,
		Payload:  []byte("Ff"), // Simplified format string
	}

	// 3. Parent KLV containing the children
	parentKLV := KLV{
		Children: []KLV{scenChild, typeChild},
	}

	// 4. Expected Output
	expected := []Scene{
		{Type: FourCCScene("SNOW"), Prob: 0.8},
		{Type: FourCCScene("URBA"), Prob: 0.6},
		{Type: FourCCScene("INDO"), Prob: 0.4},
	}

	// 5. Run the function
	result := extractSceneData(parentKLV)

	// 6. Assertion
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
