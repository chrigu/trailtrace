package parser

import "gopro/internal"

type Face struct {
	Confidence float32
	ID         int
	X          float32
	Y          float32
	W          float32
	H          float32
	Smile      float32
	Blink      float32
}

func ParseFaceData(klvs []KLV) [][]Face {
	return extractSensorData(klvs, "Face Coordinates and details", extractcFaceData)
}

// todo: handle tick tock data
func extractcFaceData(klv KLV) []Face {
	// struct ver,confidence %,ID,x,y,w,h,smile %, blink %
	// maybe check structure and don't parse dynamic data. Use a struct instead
	// BBSSSSSBB

	// todo: extract types dynamically
	// todo: handle repeat
	// todo: handle tick tock data
	// todo: handle multiple faces
	var format string = ""
	var payloads [][]byte = make([][]byte, 0)
	var scale [][]uint16
	var repeat uint32 = 1

	for _, child := range klv.Children {
		// log("Processing child:", child.FourCC)
		switch child.FourCC {
		case "FACE":
			internal.Log("FACE: found")
			payloads = append(payloads, child.Payload)

		case "TYPE":
			internal.Log("FACE: TYPE found")
			format = readPayload(child).(string)

		case "SCAL":
			internal.Log("FACE: SCAL found")
			scal := readPayload(child).([][]uint16)
			if len(scal) > 0 {
				scale = scal
			} else {
				internal.Log("Error: ParsedData is not of type []unit16")
			}
			internal.Log("FACE: scale:", scale)
		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	faces := make([]Face, len(payloads))
	for i, payload := range payloads {
		rawValues, err := parseDynamicStructure(payload, format, repeat) // todo get from gopro, honor repeat
		if err != nil {
			internal.Log("Error parsing dynamic structure:", err)
			continue
		}

		for _, values := range rawValues {
			if len(values) < 9 {
				internal.Log("Incomplete face data at index %d", i)
				continue
			}

			// only handle version 4
			if len(values) == 0 || int(float32(values[0].(uint8))/float32(scale[0][0])) != 4 {
				internal.Log("Error: No data found or version mismatch")
				faces = append(faces, Face{})
				continue
			}

			// Extract and convert the values with proper type assertions
			confidence, ok1 := values[1].(uint8)
			id, ok2 := values[2].(uint16)
			x, ok3 := values[3].(uint16)
			y, ok4 := values[4].(uint16)
			w, ok5 := values[5].(uint16)
			h, ok6 := values[6].(uint16)
			smile, ok7 := values[7].(uint8)
			blink, ok8 := values[8].(uint8)

			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 {
				internal.Log("Type assertion failed for face data at index %d", i)
				continue
			}

			faces[i] = Face{
				Confidence: float32(confidence) / float32(scale[1][0]),
				ID:         int(float32(id) / float32(scale[2][0])),
				X:          float32(x) / float32(scale[3][0]),
				Y:          float32(y) / float32(scale[4][0]),
				W:          float32(w) / float32(scale[5][0]),
				H:          float32(h) / float32(scale[6][0]),
				Smile:      float32(smile) / float32(scale[7][0]),
				Blink:      float32(blink) / float32(scale[8][0]),
			}
		}
	}

	internal.Log("faceRawData:", faces)

	if len(faces) == 0 {
		return []Face{}
	}

	return faces
}
