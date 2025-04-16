package parser

import (
	"reflect"
	"testing"
)

func TestExtractFaceData(t *testing.T) {
	// 1. Prepare FACE child (face data)
	// Format: BBSSSSSBB (version,confidence,id,x,y,w,h,smile,blink)
	faceValues := [][]interface{}{
		{uint8(4), uint8(100), uint16(0), uint16(0), uint16(0), uint16(0), uint16(0), uint8(50), uint8(0), uint16(0)}, // [version,confidence,id,x,y,w,h,smile,blink]
		{uint8(4), uint8(80), uint16(1), uint16(100), uint16(100), uint16(200), uint16(200), uint8(75), uint8(25), uint16(0)},
	}

	// Create a payload from the face values
	faceChilds := []KLV{}
	for _, packet := range faceValues {
		var facePayload []byte
		for _, value := range packet {
			switch v := value.(type) {
			case uint8:
				facePayload = append(facePayload, v)
			case uint16:
				facePayload = append(facePayload, byte(v>>8), byte(v))
			}
		}
		faceChilds = append(faceChilds, KLV{
			FourCC:   "FACE",
			DataType: int('?'),
			Repeat:   uint32(1),
			DataSize: 14,
			Payload:  facePayload,
		})
	}

	// 2. Create TYPE child with format string
	typeChild := KLV{
		FourCC:   "TYPE",
		DataType: int('c'),            // char
		DataSize: 9,                   // "BBSSSSSBB" format string
		Payload:  []byte("BBSSSSSBB"), // 000 is the padding
	}

	// 3. Create SCAL child with scale values
	// Each scale value is a uint16 (2 bytes)
	// We need 9 scale values: [1, 1, 1, 1, 1, 1, 1, 1, 1]
	scalePayload := make([]byte, 18) // 9 values * 2 bytes each
	for i := 0; i < 9; i++ {
		// Set each scale value to 1 (0x0001 in big-endian)
		scalePayload[i*2] = 0
		scalePayload[i*2+1] = 1
	}

	scaleChild := KLV{
		FourCC:   "SCAL",
		DataType: int('S'), // uint16_t
		Repeat:   9,        // 9 scale values
		DataSize: 2,        // 2 bytes per scale value
		Payload:  scalePayload,
	}

	// 4. Parent KLV containing the FACE, TYPE, and SCAL children
	parentKLV := KLV{
		Children: []KLV{
			typeChild,
			scaleChild,
			faceChilds[0],
			faceChilds[1],
		},
	}

	// 5. Expected Output
	expected := []Face{
		{
			Confidence: 100.0,
			ID:         0,
			X:          0.0,
			Y:          0.0,
			W:          0.0,
			H:          0.0,
			Smile:      50.0,
			Blink:      0.0,
		},
		{
			Confidence: 80.0,
			ID:         1,
			X:          100.0,
			Y:          100.0,
			W:          200.0,
			H:          200.0,
			Smile:      75.0,
			Blink:      25.0,
		},
	}

	// 6. Run the function
	result := extractcFaceData(parentKLV)

	// 7. Assertion
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
