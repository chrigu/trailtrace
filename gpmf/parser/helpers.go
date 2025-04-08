package parser

import (
	"fmt"
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
