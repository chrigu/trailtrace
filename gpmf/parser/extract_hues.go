package parser

import (
	"gopro/internal"
)

type Hue struct {
	Hue    uint8
	Weight uint8
}

// todo: rename to hue
func ParseHueData(klvs []KLV) [][]Hue {
	return extractSensorData(klvs,
		"Predominant hue[[hue, weight], ...]",
		extractHueData)
}

func extractHueData(klv KLV) []Hue {

	// todo: extract types dynamically
	// todo: extract types dynamically
	var format string = ""
	var payload []byte = make([]byte, 0)
	var repeat uint32 = 0

	for _, child := range klv.Children {

		switch child.FourCC {
		case "HUES":
			payload = child.Payload
			repeat = child.Repeat

		case "TYPE":
			format = readPayload(child).(string)

		default:
			continue
		}
	}

	hueRaw, err := parseDynamicStructure(payload, format, repeat) // make easier, check type and make struct
	if err != nil {
		internal.Log("Error parsing dynamic structure:", err)
		return []Hue{}
	}

	// Each color has 3 hues (RGB), so we need to process them in groups of 3
	hues := make([]Hue, len(hueRaw))

	for i := 0; i < len(hueRaw); i += 3 {
		if i+2 >= len(hueRaw) {
			break // Not enough data for a complete color
		}

		// Extract the hue values for each color component
		colorHue1, ok1 := hueRaw[i][0].(uint8)
		colorWeight1, ok2 := hueRaw[i][1].(uint8)
		colorHue2, ok3 := hueRaw[i+1][0].(uint8)
		colorWeight2, ok4 := hueRaw[i+1][1].(uint8)
		colorHue3, ok5 := hueRaw[i+2][0].(uint8)
		colorWeight3, ok6 := hueRaw[i+2][1].(uint8)

		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
			internal.Log("Type assertion failed for color at index %d", i/3)
			continue
		}

		// Add each hue to the slice
		hues[i] = Hue{Hue: colorHue1, Weight: colorWeight1}
		hues[i+1] = Hue{Hue: colorHue2, Weight: colorWeight2}
		hues[i+2] = Hue{Hue: colorHue3, Weight: colorWeight3}
	}

	return hues
}
