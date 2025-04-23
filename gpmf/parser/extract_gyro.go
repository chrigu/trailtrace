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

	var payload [][]int16
	var scale []int16

	for _, child := range klv.Children {

		switch child.FourCC {
		case "GYRO":
			payload = readPayload(child).([][]int16)

		case "SCAL":
			extractedScale, err := extractScale[int16](child)
			if err != nil {
				return []Gyroscope{}
			}
			if s, ok := extractedScale.([]int16); ok {
				scale = s
			} else {
				internal.Log("Error: ParsedData is not of type []int16")
			}
		default:
			continue
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
