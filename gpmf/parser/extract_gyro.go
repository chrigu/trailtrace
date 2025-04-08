package parser

import (
	"gopro/internal"
)

// Rename or add Accelerometer struct
type Gyroscope struct {
	X float32
	Y float32
	Z float32
}

func ParseGyroscopeData(klvs []KLV) [][]Gyroscope {
	return extractSensorData(klvs, "Gyroscope", extractGyroscopeData)
}

func extractGyroscopeData(klv KLV) []Gyroscope {
	// log("Processing STRM children", len(klv.Children))

	var payload [][]int16
	var scale []int16

	for _, child := range klv.Children {
		// log("Processing child:", child.FourCC)

		switch child.FourCC {
		case "GYRO":
			//log("GYRO found")
			payload = readPayload(child).([][]int16)

		case "SCAL":
			//log("SCAL found")
			scal := readPayload(child).([][]int16)
			if len(scal[0]) > 0 {
				scale = scal[0]
			} else {
				internal.Log("Error: ParsedData is not of type []int32")
			}
		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	gyroData := make([]Gyroscope, len(payload))
	for i := range payload {
		gyroData[i] = Gyroscope{
			X: float32(payload[i][0]) / float32(scale[0]),
			Y: float32(payload[i][1]) / float32(scale[0]),
			Z: float32(payload[i][2]) / float32(scale[0]),
		}
	}

	return gyroData

}
