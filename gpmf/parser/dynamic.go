package parser

import (
	"encoding/binary"
	"fmt"
	"math"

	"gopro/internal"
)

// parseDynamicStructure dynamically parses a buffer based on the format string
func parseDynamicStructure(data []byte, format string) ([]interface{}, error) {
	internal.Log("Parsing dynamic structure with format:", format)
	offset := 0
	totalSize := len(data)

	if totalSize == 0 {
		internal.Log("Error: No data to parse")
		return []interface{}{}, nil
	}

	values := []interface{}{} // Slice to store parsed values

	for i, char := range format {
		switch char {
		case 'B': // 8-bit unsigned integer
			if offset > totalSize {
				internal.Log("Error: Not enough data for int6 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for int8 at position %d", i)
			}
			value := data[offset]
			// log("l[%d]: %d (int32)\n", i, value)
			values = append(values, value)
			offset += 1
		case 'l': // 32-bit signed integer
			if offset+4 > totalSize {
				internal.Log("Error: Not enough data for int32 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for int32 at position %d", i)
			}
			value := int32(binary.BigEndian.Uint32(data[offset : offset+4]))
			// log("l[%d]: %d (int32)\n", i, value)
			values = append(values, value)
			offset += 4

		case 'S': // 16-bit unsigned integer
			if offset+2 > totalSize {
				internal.Log("Error: Not enough data for uint16 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for uint16 at position %d", i)
			}
			value := binary.BigEndian.Uint16(data[offset : offset+2])
			// log("S[%d]: %d (uint16)\n", i, value)
			values = append(values, value)
			offset += 2

		case 'f': // 32-bit float
			if offset+4 > totalSize {
				internal.Log("Error: Not enough data for float32 at position %d\n", i)
				return nil, fmt.Errorf("Not enough data for float32 at position %d", i)
			}
			value := math.Float32frombits(binary.BigEndian.Uint32(data[offset : offset+4]))
			// log("f[%d]: %f (float32)\n", i, value)
			values = append(values, value)
			offset += 4

		default:
			internal.Log("Unknown format character: %c\n", char)
			return nil, fmt.Errorf("Unknown format character: %c", char)
		}
	}

	// Calculate padding
	padding := (4 - (offset % 4)) % 4
	if padding > 0 && offset+int(padding) <= totalSize {
		internal.Log("Padding bytes: %d\n", padding)
		offset += int(padding)
	}

	internal.Log("Total bytes processed: %d\n", offset)
	return values, nil
}
