package parser

import (
	"slices"
	"strings"
)

// extractSensorData is a generic function to extract sensor data from KLVs
func extractSensorData[T any](klvs []KLV, sensorType string, extractFunc func(KLV) []T) [][]T {
	var dataList [][]T

	for _, klv := range klvs {
		if klv.FourCC == "STRM" {
			index := slices.IndexFunc(klv.Children, func(child KLV) bool {
				return strings.TrimSpace(string(child.Payload)) == sensorType
			})
			if index != -1 {
				data := extractFunc(klv)
				dataList = append(dataList, data)
			}
		}
		// Recursively check children
		if len(klv.Children) > 0 {
			childData := extractSensorData(klv.Children, sensorType, extractFunc)
			dataList = append(dataList, childData...)
		}
	}
	return dataList
}

func ParseGPS9Data(klvs []KLV) [][]GPS9 {
	return extractSensorData(klvs,
		"GPS (Lat., Long., Alt., 2D, 3D, days, secs, DOP, fix)",
		extractGpsData)
}

func ParseGyroscopeData(klvs []KLV) [][]Gyroscope {
	return extractSensorData(klvs, "Gyroscope", extractGyroscopeData)
}

func ParseAccelerometerData(klvs []KLV) [][]Gyroscope {
	return extractSensorData(klvs, "Accelerometer", extractAccelerometerData)
}

func ParseFaceData(klvs []KLV) [][]Face {
	return extractSensorData(klvs, "Face Coordinates and details", extractcFaceData)
}

func ParseHueData(klvs []KLV) [][]Face {
	return extractSensorData(klvs, "Face Coordinates and details", extractcFaceData)
}
