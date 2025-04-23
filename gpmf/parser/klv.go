package parser

import (
	"fmt"

	"gopro/internal"
)

type KLV struct {
	FourCC   string
	DataType int
	DataSize uint32
	Repeat   uint32
	Payload  []byte
	Children []KLV
}

func ParseGPMF(data []byte) []KLV {
	var offset uint32 = 0
	var klvs []KLV = make([]KLV, 0)

	for {
		newOffset := readKLV(data, offset, &klvs)

		if newOffset <= offset { // Stops infinite loop when offset is not advancing
			fmt.Println("Error: Offset did not advance, stopping.")
			break
		}

		offset = newOffset

		if offset >= uint32(len(data)) { // Ensures we don't read beyond available data
			break
		}
	}

	return klvs
}

func readKLV(data []byte, offset uint32, klvs *[]KLV) uint32 {
	// Check if enough bytes exist before reading

	dataOffset := offset + 8

	if dataOffset > uint32(len(data)) {
		fmt.Println("Error: Not enough data for KLV header")
		return offset + 8
	}

	klv := KLV{
		FourCC:   string(data[offset : offset+4]),
		DataType: int(data[offset+4]),
		DataSize: uint32(data[offset+5]),
		Repeat:   uint32(data[offset+6])<<8 | uint32(data[offset+7]),
		Children: make([]KLV, 0),
	}

	// Ensure payload does not exceed data slice
	if dataOffset+klv.DataSize*klv.Repeat > uint32(len(data)) {
		internal.Log("Error: Payload exceeds available data")
		return dataOffset
	}

	klv.Payload = data[dataOffset : dataOffset+klv.DataSize*klv.Repeat]
	*klvs = append(*klvs, klv)

	totalSize := klv.DataSize * klv.Repeat
	padding := (4 - (totalSize % 4)) % 4

	// Process nested KLV structures
	if klv.DataType == 0 {
		nestedOffset := uint32(0) + padding

		// Process multiple nested KLVs inside the payload
		for nestedOffset < uint32(len(klv.Payload)) {
			var nestedKLVs []KLV
			nestedOffset = readKLV(klv.Payload, nestedOffset, &nestedKLVs)

			if len(nestedKLVs) > 0 {
				(*klvs)[len(*klvs)-1].Children = append((*klvs)[len(*klvs)-1].Children, nestedKLVs...)
			}
		}
	}

	return dataOffset + klv.DataSize*klv.Repeat + padding
}
