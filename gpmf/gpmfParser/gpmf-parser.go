package gpmfParser

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strings"
)

type GPS9 struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
}

// Rename or add Accelerometer struct
type Gyroscope struct {
	X float32
	Y float32
	Z float32
}

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

func parseGPS9Data(klvs []KLV) [][]GPS9 {
	return extractSensorData(klvs,
		"GPS (Lat., Long., Alt., 2D, 3D, days, secs, DOP, fix)",
		extractGpsData)
}

func parseGyroscopeData(klvs []KLV) [][]Gyroscope {
	return extractSensorData(klvs, "Gyroscope", extractGyroscopeData)
}

func parseAccelerometerData(klvs []KLV) [][]Gyroscope {
	return extractSensorData(klvs, "Accelerometer", extractAccelerometerData)
}

func parsecFaceData(klvs []KLV) [][]Face {
	return extractSensorData(klvs, "Face Coordinates and details", extractcFaceData)
}

// parseDynamicStructure dynamically parses a buffer based on the format string
func parseDynamicStructure(data []byte, format string) ([]interface{}, error) {
	log("Parsing dynamic structure with format:", format)
	offset := 0
	totalSize := len(data)

	if totalSize == 0 {
		log("Error: No data to parse")
		return []interface{}{}, nil
	}

	values := []interface{}{} // Slice to store parsed values

	for i, char := range format {
		switch char {
		case 'B': // 8-bit unsigned integer
			if offset > totalSize {
				log("Error: Not enough data for int6 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for int8 at position %d", i)
			}
			value := data[offset]
			// log("l[%d]: %d (int32)\n", i, value)
			values = append(values, value)
			offset += 1
		case 'l': // 32-bit signed integer
			if offset+4 > totalSize {
				log("Error: Not enough data for int32 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for int32 at position %d", i)
			}
			value := int32(binary.BigEndian.Uint32(data[offset : offset+4]))
			// log("l[%d]: %d (int32)\n", i, value)
			values = append(values, value)
			offset += 4

		case 'S': // 16-bit unsigned integer
			if offset+2 > totalSize {
				log("Error: Not enough data for uint16 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for uint16 at position %d", i)
			}
			value := binary.BigEndian.Uint16(data[offset : offset+2])
			// log("S[%d]: %d (uint16)\n", i, value)
			values = append(values, value)
			offset += 2

		case 'f': // 32-bit float
			if offset+4 > totalSize {
				log("Error: Not enough data for float32 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for float32 at position %d", i)
			}
			value := math.Float32frombits(binary.BigEndian.Uint32(data[offset : offset+4]))
			// log("f[%d]: %f (float32)\n", i, value)
			values = append(values, value)
			offset += 4

		default:
			log("Unknown format character: %c\n", char)
			return nil, fmt.Errorf("Unknown format character: %c", char)
		}
	}

	// Calculate padding
	padding := (4 - (offset % 4)) % 4
	if padding > 0 && offset+int(padding) <= totalSize {
		log("Padding bytes: %d\n", padding)
		offset += int(padding)
	}

	log("Total bytes processed: %d\n", offset)
	return values, nil
}

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
			log("GPS9 found")
			payload = child.Payload

		case "TYPE":
			log("TYPE found")
			format = readPayload(child).(string)

		case "SCAL":
			log("SCAL found")
			scal := readPayload(child).([][]int32)
			if len(scal) > 0 {
				scale = scal
			} else {
				log("Error: ParsedData is not of type []int32")
			}
		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	gpsRawData, err := parseDynamicStructure(payload, format) // todo get from gopro, honor repeat
	if err != nil {
		log("Error parsing dynamic structure:", err)
	}

	return []GPS9{
		{
			Latitude:  float32(gpsRawData[0].(int32)) / float32(scale[0][0]),
			Longitude: float32(gpsRawData[1].(int32)) / float32(scale[1][0]),
			Altitude:  float32(gpsRawData[2].(int32)) / float32(scale[2][0]),
		},
	}
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
				log("Error: ParsedData is not of type []int32")
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
				log("Error: ParsedData is not of type []int32")
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
			log("FACE: found")
			payloads = append(payloads, child.Payload)

		case "TYPE":
			log("FACE: TYPE found")
			format = readPayload(child).(string)

		case "SCAL":
			log("FACE: SCAL found")
			scal := readPayload(child).([][]uint16)
			if len(scal) > 0 {
				scale = scal
			} else {
				log("Error: ParsedData is not of type []unit16")
			}
			log("FACE: scale:", scale)
		default:
			//log("Unknown FourCC", klv.FourCC)
		}
	}

	faceRawData := make([]Face, 0)
	for _, payload := range payloads {
		rawValues, err := parseDynamicStructure(payload, format) // todo get from gopro, honor repeat
		if err != nil {
			log("Error parsing dynamic structure:", err)
			continue
		}

		// only handle version 4
		if len(rawValues) == 0 || int(float32(rawValues[0].(uint8))/float32(scale[0][0])) != 4 {
			log("Error: No data found or version mismatch")
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

	log("faceRawData:", faceRawData)

	if len(faceRawData) == 0 {
		return []Face{}
	}

	return faceRawData
}

func readPayload(klv KLV) any {
	switch klv.DataType {

	// case int('b'): // int8_t
	// 	log("Type: int8_t")
	case int('B'): // uint8_t
		// 	log("Type: uint8_t")
		payload := make([][]uint8, klv.Repeat)
		for i := range klv.Repeat {
			dataPackets := make([]uint8, klv.DataSize/2)
			for j := range klv.DataSize / 2 {
				offset := (i*klv.DataSize/2 + j) * 2
				dataPackets[j] = klv.Payload[offset]
			}
			payload[i] = dataPackets
		}
		return payload
	case int('c'): // ASCII character string
		//log("Type: ASCII character string")
		// use repeat
		log("Payload:", string(klv.Payload))
		return string(klv.Payload)
	// case int('d'): // double
	// 	log("Type: double (64-bit float)")
	// case int('f'): // float
	// 	log("Type: float (32-bit float)")
	// case int('F'): // FourCC
	// 	log("Type: FourCC (32-bit character key)")
	// case int('G'): // UUID
	// 	log("Type: UUID (128-bit identifier)")
	// case int('j'): // int64_t
	// 	log("Type: int64_t (64-bit signed integer)")
	// case int('J'): // uint64_t
	// 	log("Type: uint64_t (64-bit unsigned integer)")
	case int('l'): // int32_t
		//log("Type: int32_t (32-bit signed integer)")
		payload := make([][]int32, klv.Repeat)
		for i := range klv.Repeat {
			dataPackets := make([]int32, klv.DataSize/4)
			for j := range klv.DataSize / 4 {
				offset := (i*klv.DataSize/4 + j) * 4
				dataPackets[j] = int32(binary.BigEndian.Uint32(klv.Payload[offset : offset+4]))
			}
			payload[i] = dataPackets
		}
		return payload
		// (*klvs)[len(*klvs)-1].ParsedData = []any{scal}
	// case int('L'): // uint32_t
	// 	log("Type: uint32_t (32-bit unsigned integer)")

	// case int('q'): // Q15.16
	// 	log("Type: Q15.16 (fixed-point 32-bit number)")
	// case int('Q'): // Q31.32
	// 	log("Type: Q31.32 (fixed-point 64-bit number)")
	case int('s'): // int16_t
		//log("Type: int16_t (16-bit signed integer)")
		payload := make([][]int16, klv.Repeat)
		for i := range klv.Repeat {
			dataPackets := make([]int16, klv.DataSize/2)
			for j := range klv.DataSize / 2 {
				offset := (i*klv.DataSize/2 + j) * 2
				dataPackets[j] = int16(binary.BigEndian.Uint16(klv.Payload[offset : offset+2]))
			}
			payload[i] = dataPackets
		}
		return payload
	case int('S'): // uint16_t
		// log("Type: uint16_t (16-bit unsigned integer)")
		payload := make([][]uint16, klv.Repeat)
		for i := range klv.Repeat {
			dataPackets := make([]uint16, klv.DataSize/2)
			for j := range klv.DataSize / 2 {
				offset := (i*klv.DataSize/2 + j) * 2
				dataPackets[j] = binary.BigEndian.Uint16(klv.Payload[offset : offset+2])
			}
			payload[i] = dataPackets
		}
		return payload
	// case int('U'): // UTC Date and Time string
	// 	log("Type: UTC Date and Time string")
	// case int('?'): // Complex structure
	// 	log("Type: Complex structure")
	default:
		log("Unknown data type")
		return nil
	}
}
