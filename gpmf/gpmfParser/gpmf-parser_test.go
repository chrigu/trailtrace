package gpmfParser

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestExtractGyroscopeData(t *testing.T) {
	// 1. Prepare SCAL child (scale factor)
	scaleData := []int16{1000} // example scale factor
	scalePayload := make([]byte, len(scaleData)*2)
	for i, val := range scaleData {
		binary.BigEndian.PutUint16(scalePayload[i*2:], uint16(val))
	}
	scalChild := KLV{
		FourCC:   "SCAL",
		DataType: int('s'), // int32_t
		Repeat:   1,
		DataSize: 2,
		Payload:  scalePayload,
	}

	// 2. Prepare GYRO child (gyroscope raw data)
	gyroValues := [][]int16{
		{1000, 2000, 3000},
		{4000, 5000, 6000},
	}
	var gyroPayload []byte
	for _, packet := range gyroValues {
		for _, val := range packet {
			temp := make([]byte, 2)
			binary.BigEndian.PutUint16(temp, uint16(val))
			gyroPayload = append(gyroPayload, temp...)
		}
	}

	gyroChild := KLV{
		FourCC:   "GYRO",
		DataType: int('s'), // int16_t
		Repeat:   uint32(len(gyroValues)),
		DataSize: 6, // 3 * 2 bytes
		Payload:  gyroPayload,
	}

	// 3. Parent KLV containing both children
	parentKLV := KLV{
		Children: []KLV{gyroChild, scalChild},
	}

	// 4. Expected Output
	expected := []Gyroscope{
		{X: 1.0, Y: 2.0, Z: 3.0},
		{X: 4.0, Y: 5.0, Z: 6.0},
	}

	// 5. Run the function
	result := extractGyroscopeData(parentKLV)

	// 6. Assertion
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
