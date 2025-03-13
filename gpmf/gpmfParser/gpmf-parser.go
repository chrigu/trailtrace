package gpmfParser

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/fatih/color"
)

type GPS9 struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
}

func extractGPS9Data(klvs []KLV) []GPS9 {
	var gpsDataList []GPS9

	for _, klv := range klvs {
		if klv.FourCC == "STRM" {
			fmt.Println("Stream KLV found")
			index := slices.IndexFunc(klv.Children, func(child KLV) bool {
				return strings.TrimSpace(string(child.Payload)) == "GPS (Lat., Long., Alt., 2D, 3D, days, secs, DOP, fix)"
			})
			if index != -1 {
				gpsData := extractGpsData(klv)
				gpsDataList = append(gpsDataList, gpsData...)
			}

		}
		// Recursively check children
		if len(klv.Children) > 0 {
			childGPS9 := extractGPS9Data(klv.Children)
			gpsDataList = append(gpsDataList, childGPS9...)
		}

	}
	return gpsDataList
}

// parseDynamicStructure dynamically parses a buffer based on the format string
func parseDynamicStructure(data []byte, format string) ([]interface{}, error) {
	fmt.Println("Parsing dynamic structure with format:", format)
	offset := 0
	totalSize := len(data)
	values := []interface{}{} // Slice to store parsed values

	for i, char := range format {
		switch char {
		case 'l': // 32-bit signed integer
			if offset+4 > totalSize {
				fmt.Printf("Error: Not enough data for int32 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for int32 at position %d", i)
			}
			value := int32(binary.BigEndian.Uint32(data[offset : offset+4]))
			fmt.Printf("l[%d]: %d (int32)\n", i, value)
			values = append(values, value)
			offset += 4

		case 'S': // 16-bit unsigned integer
			if offset+2 > totalSize {
				fmt.Printf("Error: Not enough data for uint16 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for uint16 at position %d", i)
			}
			value := binary.BigEndian.Uint16(data[offset : offset+2])
			fmt.Printf("S[%d]: %d (uint16)\n", i, value)
			values = append(values, value)
			offset += 2

		case 'f': // 32-bit float
			if offset+4 > totalSize {
				fmt.Printf("Error: Not enough data for float32 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for float32 at position %d", i)
			}
			value := math.Float32frombits(binary.BigEndian.Uint32(data[offset : offset+4]))
			fmt.Printf("f[%d]: %f (float32)\n", i, value)
			values = append(values, value)
			offset += 4

		default:
			fmt.Printf("Unknown format character: %c\n", char)
			return nil, fmt.Errorf("Unknown format character: %c", char)
		}
	}

	// Calculate padding
	padding := (4 - (offset % 4)) % 4
	if padding > 0 && offset+int(padding) <= totalSize {
		fmt.Printf("Padding bytes: %d\n", padding)
		offset += int(padding)
	}

	fmt.Printf("Total bytes processed: %d\n", offset)
	return values, nil
}

func extractGpsData(klv KLV) []GPS9 {
	fmt.Println("Processing STRM children", len(klv.Children))

	var format string = ""
	var payload []byte = make([]byte, 0)
	var scale []int32

	for _, child := range klv.Children {
		fmt.Println("Processing child:", child.FourCC)

		switch child.FourCC {
		case "GPS9":
			color.Green("GPS9 found")
			payload = child.Payload

		case "TYPE":
			fmt.Println("TYPE found")
			format = readPayload(child).(string)

		case "SCAL":
			fmt.Println("SCAL found")
			scal := readPayload(child).([]int32)
			if len(scal) > 0 {
				scale = scal
			} else {
				fmt.Println("Error: ParsedData is not of type []int32")
			}
		default:
			//fmt.Println("Unknown FourCC", klv.FourCC)
		}
	}

	gpsRawData, err := parseDynamicStructure(payload, format) // todo get from gopro
	if err != nil {
		fmt.Println("Error parsing dynamic structure:", err)
	}

	gpsData := []GPS9{
		{
			Latitude:  float32(gpsRawData[0].(int32)) / float32(scale[0]),
			Longitude: float32(gpsRawData[1].(int32)) / float32(scale[1]),
			Altitude:  float32(gpsRawData[2].(int32)) / float32(scale[2]),
		},
	}

	return gpsData

}

func readPayload(klv KLV) any {
	switch klv.DataType {

	// case int('b'): // int8_t
	// 	fmt.Println("Type: int8_t")
	// case int('B'): // uint8_t
	// 	fmt.Println("Type: uint8_t")
	case int('c'): // ASCII character string
		fmt.Println("Type: ASCII character string")
		// use repeat
		fmt.Println("Payload:", string(klv.Payload))
		return string(klv.Payload)
	// case int('d'): // double
	// 	fmt.Println("Type: double (64-bit float)")
	// case int('f'): // float
	// 	fmt.Println("Type: float (32-bit float)")
	// case int('F'): // FourCC
	// 	fmt.Println("Type: FourCC (32-bit character key)")
	// case int('G'): // UUID
	// 	fmt.Println("Type: UUID (128-bit identifier)")
	// case int('j'): // int64_t
	// 	fmt.Println("Type: int64_t (64-bit signed integer)")
	// case int('J'): // uint64_t
	// 	fmt.Println("Type: uint64_t (64-bit unsigned integer)")
	case int('l'): // int32_t
		fmt.Println("Type: int32_t (32-bit signed integer)")
		scal := make([]int32, klv.Repeat)
		for i := range klv.Repeat {
			scal[i] = int32(binary.BigEndian.Uint32(klv.Payload[i*4 : i*4+4]))
		}
		return scal
		// (*klvs)[len(*klvs)-1].ParsedData = []any{scal}
	// case int('L'): // uint32_t
	// 	fmt.Println("Type: uint32_t (32-bit unsigned integer)")

	// case int('q'): // Q15.16
	// 	fmt.Println("Type: Q15.16 (fixed-point 32-bit number)")
	// case int('Q'): // Q31.32
	// 	fmt.Println("Type: Q31.32 (fixed-point 64-bit number)")
	// case int('s'): // int16_t
	// 	fmt.Println("Type: int16_t (16-bit signed integer)")
	// case int('S'): // uint16_t
	// 	fmt.Println("Type: uint16_t (16-bit unsigned integer)")
	// case int('U'): // UTC Date and Time string
	// 	fmt.Println("Type: UTC Date and Time string")
	// case int('?'): // Complex structure
	// 	fmt.Println("Type: Complex structure")
	default:
		fmt.Println("Unknown data type")
		return nil
	}
}
