package parser

import (
	"gopro/internal"
)

func ParseAccelerometerData(klvs []KLV) [][]Gyroscope {
	return extractSensorData(klvs, "Accelerometer", extractAccelerometerData)
}

func extractAccelerometerData(klv KLV) []Gyroscope {
	// log("Processing STRM children", len(klv.Children))

	var payload [][]int16
	var scale []int16

	for _, child := range klv.Children {
		// log("Processing child:", child.FourCC)

		switch child.FourCC {
		case "ACCL":
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
