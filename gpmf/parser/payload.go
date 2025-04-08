package parser

import (
	"encoding/binary"

	"gopro/internal"
)

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
		internal.Log("Payload:", string(klv.Payload))
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
		internal.Log("Unknown data type")
		return nil
	}
}
