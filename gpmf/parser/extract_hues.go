package parser

import (
	"gopro/internal"
)

type Hue struct {
	Hue    uint8
	Weight uint8
}

type Color struct {
	Red   Hue
	Green Hue
	Blue  Hue
}

// todo: rename to hue
func ParseColorData(klvs []KLV) [][]Color {
	return extractSensorData(klvs,
		"Predominant hue[[hue, weight], ...]",
		extractHueData)
}

func extractHueData(klv KLV) []Color {
	// log("Processing STRM children", len(klv.Children))

	// todo: extract types dynamically
	// todo: extract types dynamically
	var format string = ""
	var payload []byte = make([]byte, 0)
	var repeat uint32 = 0

	for _, child := range klv.Children {
		// log("Processing child:", child.FourCC)

		switch child.FourCC {
		case "HUES":
			internal.Log("HUES found")
			payload = child.Payload
			repeat = child.Repeat

		case "TYPE":
			internal.Log("TYPE found")
			format = readPayload(child).(string)

		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	hueRaw, err := parseDynamicStructure(payload, format, repeat) // make easier, check type and make struct
	if err != nil {
		internal.Log("Error parsing dynamic structure:", err)
		return []Color{}
	}

	// Each color has 3 hues (RGB), so we need to process them in groups of 3
	colors := make([]Color, len(hueRaw)/3)

	for i := 0; i < len(hueRaw); i += 3 {
		if i+2 >= len(hueRaw) {
			break // Not enough data for a complete color
		}

		// Extract the hue values for each color component
		redHue, ok1 := hueRaw[i][0].(uint8)
		redWeight, ok2 := hueRaw[i][1].(uint8)
		greenHue, ok3 := hueRaw[i+1][0].(uint8)
		greenWeight, ok4 := hueRaw[i+1][1].(uint8)
		blueHue, ok5 := hueRaw[i+2][0].(uint8)
		blueWeight, ok6 := hueRaw[i+2][1].(uint8)

		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
			internal.Log("Type assertion failed for color at index %d", i/3)
			continue
		}

		colors[i/3] = Color{
			Red:   Hue{Hue: redHue, Weight: redWeight},
			Green: Hue{Hue: greenHue, Weight: greenWeight},
			Blue:  Hue{Hue: blueHue, Weight: blueWeight},
		}
	}

	return colors
}
