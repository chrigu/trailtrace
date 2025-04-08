package gpmfParser

import (
	"fmt"

	"gopro/parser"
)

// PrintTree recursively prints the KLV hierarchy in a tree structure
func PrintTree(klvs []parser.KLV, prefix string) {
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

// func extractGPS9Data(klvs []KLV) []GPS9 {
// 	var gpsDataList []GPS9

// 	for _, klv := range klvs {
// 		if klv.FourCC == "STRM" {

// 			for _, data := range klv.ParsedData {
// 				if gps, ok := data.(GPS9); ok {
// 					gpsDataList = append(gpsDataList, gps)
// 				} else {
// 					fmt.Println("STRM data:", klv.ParsedData)
// 					fmt.Println("Warning: ParsedData entry is not of type GPS9")
// 				}
// 			}
// 		}

// 		// Recursively check children
// 		if len(klv.Children) > 0 {
// 			childGPS9 := extractGPS9Data(klv.Children)
// 			gpsDataList = append(gpsDataList, childGPS9...)
// 		}
// 	}

// 	return gpsDataList
// }
