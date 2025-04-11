package parser

import (
	"reflect"
	"testing"
)

func TestExtractLumaData(t *testing.T) {
	// 1. Prepare YAVG child (luminance data)
	lumaValues := [][]uint8{
		{100},
		{150},
		{200},
	}

	// Create a payload from the luma values
	var yavgPayload []byte
	for _, packet := range lumaValues {
		yavgPayload = append(yavgPayload, packet...)
	}

	yavgChild := KLV{
		FourCC:   "YAVG",
		DataType: int('B'), // uint8_t
		Repeat:   uint32(len(lumaValues)),
		DataSize: 1, // 1 byte per value
		Payload:  yavgPayload,
	}

	// 2. Parent KLV containing the YAVG child
	parentKLV := KLV{
		Children: []KLV{yavgChild},
	}

	// 3. Expected Output
	expected := []Luma{
		{Luminance: 100},
		{Luminance: 150},
		{Luminance: 200},
	}

	// 4. Run the function
	result := extractLumaData(parentKLV)

	// 5. Assertion
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
