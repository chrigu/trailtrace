package main

import (
	"fmt"
	"os"

	"github.com/abema/go-mp4"
)

const (
	GoProMetaName = "GoPro MET"
)

func ExtractMetadataTrack(file *os.File) (*mp4.BoxInfo, error) {
	// Extract metadata track from the MP4 file
	var metadataTrack *mp4.BoxInfo

	trackBoxes, err := mp4.ExtractBox(file, nil, mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeTrak()})

	if err != nil {
		return nil, fmt.Errorf("No tracks found: %w", err)
	}

	for _, trackBox := range trackBoxes {
		hdlrBoxes, err := mp4.ExtractBoxWithPayload(file, trackBox, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeHdlr()})
		if err != nil {
			return nil, fmt.Errorf("No hdlr boxes found: %w", err)
		}
		for _, hdlrBox := range hdlrBoxes {
			hdlr := hdlrBox.Payload.(*mp4.Hdlr)

			if string(hdlr.Name) == GoProMetaName {
				metadataTrack = trackBox
				break
			}

			if metadataTrack != nil {
				break
			}
		}
	}

	return metadataTrack, nil
}