package parser

import (
	"reflect"
	"testing"
)

func TestExtractHueData(t *testing.T) {
	// 1. Prepare HUES child (hue data)
	hueValues := [][]uint8{
		{10, 20}, // [hue, weight]
		{30, 40},
		{50, 60},
	}

	// Create a payload from the hue values
	var huesPayload []byte
	for _, packet := range hueValues {
		huesPayload = append(huesPayload, packet...)
	}

	huesChild := KLV{
		FourCC:   "HUES",
		DataType: int('B'), // uint8_t
		Repeat:   uint32(len(hueValues)),
		DataSize: 2, // 2 bytes per value (hue and weight)
		Payload:  huesPayload,
	}

	// 2. Create TYPE child with format string
	typeChild := KLV{
		FourCC:   "TYPE",
		DataType: int('c'), // char
		DataSize: 2,        // "BB" format string
		Payload:  []byte("BB"),
	}

	// 3. Parent KLV containing the HUES and TYPE children
	parentKLV := KLV{
		Children: []KLV{huesChild, typeChild},
	}

	// 4. Expected Output
	expected := []Hue{
		{Hue: 10, Weight: 20},
		{Hue: 30, Weight: 40},
		{Hue: 50, Weight: 60},
	}

	// 5. Run the function
	result := extractHueData(parentKLV)

	// 6. Assertion
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestParseHueData(t *testing.T) {
	// 1. Prepare multiple KLVs with HUES data
	klvs := []KLV{
		{
			FourCC: "STRM",
			Children: []KLV{
				{
					FourCC: "HUES",
					Children: []KLV{
						{
							FourCC:   "HUES",
							DataType: int('B'),
							Repeat:   3,
							DataSize: 2,
							Payload:  []byte{10, 20, 30, 40, 50, 60},
						},
						{
							FourCC:   "TYPE",
							DataType: int('c'),
							DataSize: 2,
							Payload:  []byte("BB"),
						},
					},
				},
				{
					FourCC:   "STNM",
					DataType: int('c'),
					DataSize: 2,
					Payload:  []byte("Predominant hue[[hue, weight], ...]"),
				},
			},
		},
	}

	// 2. Expected Output
	expected := [][]Hue{
		{
			{Hue: 10, Weight: 20},
			{Hue: 30, Weight: 40},
			{Hue: 50, Weight: 60},
		},
	}

	// 3. Run the function
	result := ParseHueData(klvs)

	// 4. Assertion
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
