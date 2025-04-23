package parser

import (
	"fmt"
	"gopro/internal"
	"slices"
	"strings"
)

// PrintTree recursively prints the KLV hierarchy in a tree structure
func PrintTree(klvs []KLV, prefix string) {
	for i, klv := range klvs {
		// Determine tree branching characters
		connector := "├──"
		if i == len(klvs)-1 {
			connector = "└──"
		}

		// Print the current node
		fmt.Printf("%s%s %s [%d, %d]\n", prefix, connector, klv.FourCC, klv.DataSize, klv.Repeat)

		// Recursively print children with adjusted prefix
		newPrefix := prefix + "│   "
		if i == len(klvs)-1 {
			newPrefix = prefix + "    "
		}
		PrintTree(klv.Children, newPrefix)
	}
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

func extractScale[T int16 | int32 | uint16](child KLV) (any, error) {
	scal, ok := readPayload(child).([][]T)
	if !ok {
		err := fmt.Errorf("failed to parse SCAL payload as [][]%T", *new(T))
		internal.Log("Error: %v", err)
		return nil, err
	}
	if len(scal) == 0 {
		err := fmt.Errorf("no scale found")
		internal.Log("Error: %v", err)
		return nil, err
	}
	// Return either []T or [][]T depending on input structure
	if len(scal) == 1 {
		return scal[0], nil
	}
	return scal, nil
}
