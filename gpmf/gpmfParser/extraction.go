package gpmfParser

import (
	"fmt"
	"io"

	"github.com/abema/go-mp4"
)

const (
	GoProMetaName = "GoPro MET"
)

type TelemetryMetadata struct {
	ChunkOffsets   []uint32
	ChunkSizes     []uint32
	SampleToChunks []mp4.StscEntry
	TimeToSamples  []mp4.SttsEntry
}

func ExtractTelemetryData(file io.ReadSeeker) []GPS9 {
	var metadataTrack *mp4.BoxInfo
	var err error

	telemetryMetadata := TelemetryMetadata{}
	// Extract metadata track from the MP4 file
	metadataTrack, err = extractMetadataTrack(file)
	if err != nil {
		fmt.Println("Error extracting metadata track:", err)
		return []GPS9{}
	}

	if metadataTrack == nil {
		fmt.Println("No metadata track found")
		return []GPS9{}
	}
	fmt.Println("metadata track", metadataTrack)

	stcoBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStco()})
	for _, stcoBox := range stcoBoxes {
		stcoBox := stcoBox.Payload.(*mp4.Stco)
		telemetryMetadata.ChunkOffsets = stcoBox.ChunkOffset
	}

	stszBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsz()})
	for _, stszBox := range stszBoxes {
		stszBox := stszBox.Payload.(*mp4.Stsz)
		telemetryMetadata.ChunkSizes = stszBox.EntrySize
	}

	// get Stsc
	stscBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsc()})
	for _, stscBox := range stscBoxes {
		stscBox := stscBox.Payload.(*mp4.Stsc)
		telemetryMetadata.SampleToChunks = stscBox.Entries
	}

	// get Stts
	sttsBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStts()})
	for _, sttsBox := range sttsBoxes {
		sttsBox := sttsBox.Payload.(*mp4.Stts)
		telemetryMetadata.TimeToSamples = sttsBox.Entries
	}

	// // Read mdat size
	// mdatBoxes, err := mp4.ExtractBox(file, nil, mp4.BoxPath{mp4.BoxTypeMdat()})
	// //error handling
	// fmt.Println("Offset mdat", mdatBoxes[0].Offset, "Size mdat", mdatBoxes[0].Size)
	// telemetryMetadata.Offset = mdatBoxes[0].Offset

	fmt.Println("Telemetry Metadata", telemetryMetadata.ChunkOffsets[0])

	data, _ := readRawData(file, &telemetryMetadata)
	// writeBinaryToFile("telemetry.bin", data)

	// move elsewhere
	klvs := ParseGPMF(data)
	gpsData := extractGPS9Data(klvs)
	fmt.Println("GPS9 data:", len(gpsData))
	// for _, gps := range gpsData {
	// 	fmt.Println("GPS9 data:", gps)
	// }

	if err != nil {
		fmt.Println("Error reading MP4 structure:", err)
	}

	return gpsData
}

func extractMetadataTrack(file io.ReadSeeker) (*mp4.BoxInfo, error) {
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

func readRawData(file io.ReadSeeker, telemetryMetadata *TelemetryMetadata) ([]byte, error) {
	var totalSize uint32
	for _, chunkSize := range telemetryMetadata.ChunkSizes {
		totalSize += chunkSize
	}
	buffer := make([]byte, totalSize)

	var bufferPos uint64 = 0
	for idx, offset := range telemetryMetadata.ChunkOffsets {
		chunkSize := uint64(telemetryMetadata.ChunkSizes[idx])

		_, err := file.Seek(int64(offset), io.SeekStart)
		if err != nil {
			fmt.Printf("Error seeking at offset %d: %v\n", offset, err)
			return nil, err
		}

		n, err := file.Read(buffer[bufferPos : bufferPos+chunkSize])
		if err != nil {
			fmt.Printf("Error reading at offset %d: %v\n", offset, err)
			return nil, err
		}
		fmt.Printf("Read %d bytes from offset %d\n", n, offset)
		bufferPos += chunkSize
	}
	return buffer, nil
}
