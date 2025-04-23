package mp4

import (
	"fmt"
	"io"
	"strings"

	"gopro/internal"

	"github.com/abema/go-mp4"
)

// todo: Maybe add a list of klvs to return
func ExtractTelemetryFromMp4(file io.ReadSeeker) ([]byte, TelemetryMetadata) {
	var metadataTrack *mp4.BoxInfo
	var err error

	telemetryMetadata := TelemetryMetadata{}
	// Extract metadata track from the MP4 file
	metadataTrack, err = extractMetadataTrack(file)
	if err != nil {
		internal.Log("Error extracting metadata track:", err)
		return nil, telemetryMetadata
	}

	if metadataTrack == nil {
		internal.Log("No metadata track found")
		return nil, telemetryMetadata
	}

	// extract Mdhd
	mdhdBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMdhd()})
	telemetryMetadata.TimeScale = mdhdBoxes[0].Payload.(*mp4.Mdhd).Timescale
	telemetryMetadata.CreationTime = convertMP4TimeToUnixMs(mdhdBoxes[0].Payload.(*mp4.Mdhd).CreationTimeV0)

	// extract Stco
	stcoBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStco()})
	if err != nil {
		internal.Log("Error extracting Stco:", err)
	}
	for _, stcoBox := range stcoBoxes {
		stcoBox := stcoBox.Payload.(*mp4.Stco)
		telemetryMetadata.ChunkOffsets = stcoBox.ChunkOffset
	}

	// extract Stsz
	stszBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsz()})
	if err != nil {
		internal.Log("Error extracting Stsz:", err)
	}
	for _, stszBox := range stszBoxes {
		stszBox := stszBox.Payload.(*mp4.Stsz)
		telemetryMetadata.ChunkSizes = stszBox.EntrySize
	}

	// extract Stsc
	stscBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsc()})
	if err != nil {
		internal.Log("Error extracting Stsc:", err)
	}
	for _, stscBox := range stscBoxes {
		stscBox := stscBox.Payload.(*mp4.Stsc)
		telemetryMetadata.SampleToChunks = stscBox.Entries
	}

	// extract Stts
	sttsBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStts()})
	if err != nil {
		internal.Log("Error extracting Stts:", err)
	}
	for _, sttsBox := range sttsBoxes {
		sttsBox := sttsBox.Payload.(*mp4.Stts)
		telemetryMetadata.TimeToSamples = sttsBox.Entries
	}

	// Read mdat size
	mdatBoxes, err := mp4.ExtractBox(file, nil, mp4.BoxPath{mp4.BoxTypeMdat()})
	if err != nil {
		internal.Log("Error extracting mdat:", err)
	}
	internal.Log("Offset mdat", mdatBoxes[0].Offset, "Size mdat", mdatBoxes[0].Size)

	data, err := readRawData(file, &telemetryMetadata)
	if err != nil {
		internal.Log("Error reading raw data:", err)
	}

	return data, telemetryMetadata
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

			if strings.TrimSpace(hdlr.Name) == GoProMetaName {
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
			internal.Log("Error seeking at offset %d: %v\n", offset, err)
			return nil, err
		}

		_, err = file.Read(buffer[bufferPos : bufferPos+chunkSize])
		if err != nil {
			internal.Log("Error reading at offset %d: %v\n", offset, err)
			return nil, err
		}
		bufferPos += chunkSize
	}
	return buffer, nil
}

func convertMP4TimeToUnixMs(creationTimeV0 uint32) int64 {
	// MP4 Epoch starts at 1904-01-01, Unix Epoch starts at 1970-01-01
	mp4EpochOffset := int64(2082844800)

	// Convert to Unix timestamp
	return (int64(creationTimeV0) - mp4EpochOffset) * 1000
}
