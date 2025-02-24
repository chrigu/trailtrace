package gpmfParser

import (
	"fmt"
	// "encoding/binary"
)

func ParseGPMF(data []byte) {
	fmt.Println("Hello from gpmf-parser")
	var offset uint32 = 0

	for {
		newOffset := readKLV(data, offset)

		if newOffset <= offset { // Stops infinite loop when offset is not advancing
			fmt.Println("Error: Offset did not advance, stopping.")
			break
		}

		fmt.Println("Offset advanced to:", newOffset)
		offset = newOffset

		// if offset >= uint32(len(data)) { // Ensures we don't read beyond available data
		if offset >= 5000 { // Ensures we don't read beyond available data
			break
		}
	}
}

func readKLV(data []byte, offset uint32) uint32 {
	// Check if enough bytes exist before reading
	if offset+8 > uint32(len(data)) {
		fmt.Println("Error: Not enough data for KLV header")
		return offset + 8
	}

	fourCC := string(data[offset : offset+4])
	dataType := int(data[offset+4])
	dataSize := uint32(data[offset+5])
	repeat := uint32(data[offset+6])<<8 | uint32(data[offset+7])

	// Ensure payload does not exceed data slice
	if offset+8+dataSize*repeat > uint32(len(data)) {
		// fmt.Println("Error: Payload exceeds available data")
		return offset + 8
	}

	payload := data[offset+8 : offset+8+dataSize*repeat]
	totalSize := dataSize * repeat
	padding := (4 - (totalSize % 4)) % 4

	fmt.Println("FourCC:", fourCC, "DataType:", dataType, "DataSize:", dataSize, "Repeat:", repeat, "Padding:", padding)

	// Process nested KLV structures

	switch dataType {
	case 0:
		fmt.Println("Processing nested KLV entries")
		nestedOffset := uint32(0) + padding

		// Process multiple nested KLVs inside the payload
		for nestedOffset < uint32(len(payload)) {
			nestedOffset = readKLV(payload, nestedOffset)
		}
	case int('b'): // int8_t
		fmt.Println("Type: int8_t")
	case int('B'): // uint8_t
		fmt.Println("Type: uint8_t")
	case int('c'): // ASCII character string
		fmt.Println("Type: ASCII character string")
		fmt.Println("Payload:", string(payload))
	case int('d'): // double
		fmt.Println("Type: double (64-bit float)")
	case int('f'): // float
		fmt.Println("Type: float (32-bit float)")
	case int('F'): // FourCC
		fmt.Println("Type: FourCC (32-bit character key)")
	case int('G'): // UUID
		fmt.Println("Type: UUID (128-bit identifier)")
	case int('j'): // int64_t
		fmt.Println("Type: int64_t (64-bit signed integer)")
	case int('J'): // uint64_t
		fmt.Println("Type: uint64_t (64-bit unsigned integer)")
	case int('l'): // int32_t
		fmt.Println("Type: int32_t (32-bit signed integer)")
	case int('L'): // uint32_t
		fmt.Println("Type: uint32_t (32-bit unsigned integer)")
	case int('q'): // Q15.16
		fmt.Println("Type: Q15.16 (fixed-point 32-bit number)")
	case int('Q'): // Q31.32
		fmt.Println("Type: Q31.32 (fixed-point 64-bit number)")
	case int('s'): // int16_t
		fmt.Println("Type: int16_t (16-bit signed integer)")
	case int('S'): // uint16_t
		fmt.Println("Type: uint16_t (16-bit unsigned integer)")
	case int('U'): // UTC Date and Time string
		fmt.Println("Type: UTC Date and Time string")
	case int('?'): // Complex structure
		fmt.Println("Type: Complex structure")
	default:
		fmt.Println("Unknown data type")
	}

	return offset + 8 + dataSize*repeat + padding
}
