package main

import (
	"fmt"
	"os"

	"github.com/abema/go-mp4"
)


func main() {
	// Open MP4 file for reading
	var metadataTrack *mp4.BoxInfo

	file, err := os.Open("../GX010025.mp4")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Extract metadata track from the MP4 file
	metadataTrack, err = ExtractMetadataTrack(file)
	if err != nil {
		fmt.Println("Error extracting metadata track:", err)
		return
	}

	if metadataTrack == nil {
		fmt.Println("No metadata track found")
		return
	}
	fmt.Println("metadata track", metadataTrack)

	stcoBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStco()})
	for _, stcoBox := range stcoBoxes {
		stcoBox := stcoBox.Payload.(*mp4.Stco)
		fmt.Println("stco", stcoBox.ChunkOffset)
	}

	stszBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsz()})
	for _, stszBox := range stszBoxes {
		stszBox := stszBox.Payload.(*mp4.Stsz)
		fmt.Println("stco", stszBox.EntrySize)
	}

	// get Stsc
	stscBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsc()})
	for _, stscBox := range stscBoxes {
		stscBox := stscBox.Payload.(*mp4.Stsc)
		fmt.Println("stsc", stscBox.Entries)
	}

	// get Stts
	sttsBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStts()})
	for _, sttsBox := range sttsBoxes {
		sttsBox := sttsBox.Payload.(*mp4.Stts)
		fmt.Println("stts", sttsBox.Entries)
	}

	// Read mdat size
	mdatBoxes, err := mp4.ExtractBox(file, nil, mp4.BoxPath{mp4.BoxTypeMdat()})
	//error handling
	fmt.Println("Offset mdat", mdatBoxes[0].Offset, "Size mdat", mdatBoxes[0].Size)

	if err != nil {
		fmt.Println("Error reading MP4 structure:", err)
	}
}

