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

// todo: handle tick tock data
func extractcFaceData(klv KLV) []Face {
	// struct ver,confidence %,ID,x,y,w,h,smile %, blink %
	// BBSSSSSBB

	// todo: extract types dynamically
	// todo: handle repeat
	// todo: handle tick tock data
	// todo: handle multiple faces
	var format string = ""
	var payloads [][]byte = make([][]byte, 0)
	var scale [][]uint16

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

	faceRawData := make([]Face, 0)
	for _, payload := range payloads {
		rawValues, err := parseDynamicStructure(payload, format) // todo get from gopro, honor repeat
		if err != nil {
			internal.Log("Error parsing dynamic structure:", err)
			continue
		}

		// only handle version 4
		if len(rawValues) == 0 || int(float32(rawValues[0].(uint8))/float32(scale[0][0])) != 4 {
			internal.Log("Error: No data found or version mismatch")
			faceRawData = append(faceRawData, Face{})
			continue
		}

		face := Face{
			Confidence: float32(rawValues[1].(uint8)) / float32(scale[1][0]),
			ID:         int(float32(rawValues[2].(uint16)) / float32(scale[2][0])),
			X:          float32(rawValues[3].(uint16)) / float32(scale[3][0]),
			Y:          float32(rawValues[4].(uint16)) / float32(scale[4][0]),
			W:          float32(rawValues[5].(uint16)) / float32(scale[5][0]),
			H:          float32(rawValues[6].(uint16)) / float32(scale[6][0]),
			Smile:      float32(rawValues[7].(uint8)) / float32(scale[7][0]),
			Blink:      float32(rawValues[8].(uint8)) / float32(scale[8][0]),
		}
		faceRawData = append(faceRawData, face)
	}

	internal.Log("faceRawData:", faceRawData)

	if len(faceRawData) == 0 {
		return []Face{}
	}

	return faceRawData
}
