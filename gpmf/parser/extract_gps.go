package parser

import (
	"gopro/internal"
)

func extractGpsData(klv KLV) []GPS9 {
	// log("Processing STRM children", len(klv.Children))

	// todo: extract types dynamically
	var format string = ""
	var payload []byte = make([]byte, 0)
	var scale [][]int32

	for _, child := range klv.Children {
		// log("Processing child:", child.FourCC)

		switch child.FourCC {
		case "GPS9":
			internal.Log("GPS9 found")
			payload = child.Payload

		case "TYPE":
			internal.Log("TYPE found")
			format = readPayload(child).(string)

		case "SCAL":
			internal.Log("SCAL found")
			scal := readPayload(child).([][]int32)
			if len(scal) > 0 {
				scale = scal
			} else {
				internal.Log("Error: ParsedData is not of type []int32")
			}
		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	gpsRawData, err := parseDynamicStructure(payload, format) // todo get from gopro, honor repeat
	if err != nil {
		internal.Log("Error parsing dynamic structure:", err)
	}

	return []GPS9{
		{
			Latitude:  float32(gpsRawData[0].(int32)) / float32(scale[0][0]),
			Longitude: float32(gpsRawData[1].(int32)) / float32(scale[1][0]),
			Altitude:  float32(gpsRawData[2].(int32)) / float32(scale[2][0]),
		},
	}
}
