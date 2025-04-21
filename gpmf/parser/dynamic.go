package parser

import (
	"encoding/binary"
	"fmt"
	"math"

	"gopro/internal"
)

// parseDynamicStructure dynamically parses a buffer based on the format string
func parseDynamicStructure(data []byte, format string, repeat uint32) ([][]interface{}, error) {
	internal.Log("Parsing dynamic structure with format:", format)
	offset := 0
	totalSize := len(data)

	if totalSize == 0 {
		internal.Log("Error: No data to parse")
		return [][]interface{}{}, nil
	}

	// Create a slice to store all repeated structures
	allValues := make([][]interface{}, 0, repeat)

	// Calculate the size of one structure
	structureSize := 0
	for _, char := range format {
		switch char {
		case 'B': // 8-bit unsigned integer
			structureSize += 1
		case 'l': // 32-bit signed integer
			structureSize += 4
		case 'S': // 16-bit unsigned integer
			structureSize += 2
		case 'f': // 32-bit float
			structureSize += 4
		}
	}

	// Add padding to structure size
	structureSize = (structureSize*int(repeat) + 3) & ^3 // Round up to nearest multiple of 4
	totalSizePAdded := (totalSize + 3) & ^3

	// Check if we have enough data for all repetitions
	if offset+structureSize > totalSizePAdded {
		internal.Log("Error: Not enough data for %d repetitions", repeat)
		return nil, fmt.Errorf("Not enough data for %d repetitions", repeat)
	}

	// Parse the structure repeat times
	for r := uint32(0); r < repeat; r++ {
		values := make([]interface{}, 0, len(format))

		for i, char := range format {
			switch char {
			case 'B': // 8-bit unsigned integer
				if offset > totalSize {
					internal.Log("Error: Not enough data for int8 at position %d\n", i)
					return nil, fmt.Errorf("Not enough data for int8 at position %d", i)
				}
				value := data[offset]
				values = append(values, interface{}(value))
				offset += 1
			case 'F': // 32-bit four character key -- FourCC
				if offset+4 > totalSize {
					internal.Log("Error: Not enough data for FourCC at position %d\n", i)
					return nil, fmt.Errorf("Not enough data for FourCC at position %d", i)
				}
				value := string(data[offset : offset+4])
				values = append(values, interface{}(value))
				offset += 4
			case 'l': // 32-bit signed integer
				if offset+4 > totalSize {
					internal.Log("Error: Not enough data for int32 at position %d\n", i)
					return nil, fmt.Errorf("Not enough data for int32 at position %d", i)
				}
				value := int32(binary.BigEndian.Uint32(data[offset : offset+4]))
				values = append(values, interface{}(value))
				offset += 4

			case 'S': // 16-bit unsigned integer
				if offset+2 > totalSize {
					internal.Log("Error: Not enough data for uint16 at position %d\n", i)
					return nil, fmt.Errorf("Not enough data for uint16 at position %d", i)
				}
				value := binary.BigEndian.Uint16(data[offset : offset+2])
				values = append(values, interface{}(value))
				offset += 2

			case 'f': // 32-bit float
				if offset+4 > totalSize {
					internal.Log("Error: Not enough data for float32 at position %d\n", i)
					return nil, fmt.Errorf("Not enough data for float32 at position %d", i)
				}
				value := math.Float32frombits(binary.BigEndian.Uint32(data[offset : offset+4]))
				values = append(values, interface{}(value))
				offset += 4

			default:
				internal.Log("Unknown format character: %c\n", char)
				return nil, fmt.Errorf("Unknown format character: %c", char)
			}
		}

		// Calculate padding
		// padding := (4 - (offset % 4)) % 4
		// if padding > 0 && offset+int(padding) <= totalSize {
		// 	internal.Log("Padding bytes: %d\n", padding)
		// 	offset += int(padding)
		// }

		// Add this structure's values to the result
		allValues = append(allValues, values)
	}

	internal.Log("Total bytes processed: %d\n", offset)
	return allValues, nil
}
