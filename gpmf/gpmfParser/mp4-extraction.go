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
	CreationTime   int64
	TimeScale      uint32
	ChunkOffsets   []uint32
	ChunkSizes     []uint32
	SampleToChunks []mp4.StscEntry
	TimeToSamples  []mp4.SttsEntry
}

type TimeSample struct {
	TimeStamp int64
}

type GPSSample struct {
	GPS9
	TimeSample
}

type GyroSample struct {
	Gyroscope
	TimeSample
}

func ExtractTelemetryDataFromMp4(file io.ReadSeeker) ([]GPSSample, []GyroSample) {
	var metadataTrack *mp4.BoxInfo
	var err error

	telemetryMetadata := TelemetryMetadata{}
	// Extract metadata track from the MP4 file
	metadataTrack, err = extractMetadataTrack(file)
	if err != nil {
		fmt.Println("Error extracting metadata track:", err)
		return []GPSSample{}, []GyroSample{}
	}

	if metadataTrack == nil {
		fmt.Println("No metadata track found")
		return []GPSSample{}, []GyroSample{}
	}

	mdhdBoxes, err := mp4.ExtractBoxWithPayload(file, metadataTrack, mp4.BoxPath{mp4.BoxTypeMdia(), mp4.BoxTypeMdhd()})
	telemetryMetadata.TimeScale = mdhdBoxes[0].Payload.(*mp4.Mdhd).Timescale
	telemetryMetadata.CreationTime = getUnixTimestamp(mdhdBoxes[0].Payload.(*mp4.Mdhd).CreationTimeV0)

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

	// Read mdat size
	mdatBoxes, err := mp4.ExtractBox(file, nil, mp4.BoxPath{mp4.BoxTypeMdat()})
	//error handling
	fmt.Println("Offset mdat", mdatBoxes[0].Offset, "Size mdat", mdatBoxes[0].Size)

	fmt.Println("Telemetry Metadata", telemetryMetadata.TimeScale, telemetryMetadata.TimeToSamples)

	data, _ := readRawData(file, &telemetryMetadata)

	return ExtractTelemetryData(data, &telemetryMetadata, false)

}

func ExtractTelemetryData(data []byte, telemetryMetadata *TelemetryMetadata, printTree bool) ([]GPSSample, []GyroSample) {
	klvs := ParseGPMF(data)

	if printTree {
		PrintTree(klvs, "")
	}

	fmt.Println("KLVs", len(klvs))
	gpsData := extractGPS9Data(klvs)
	gyroData := extractGyroData(klvs)
	fmt.Println("GPS9 data:", len(gpsData), "Gyro data:", len(gyroData))
	gpsDataSamples := assignTimestampsToGps(gpsData, telemetryMetadata)
	gyroDataSamples := assignTimestampsToGyroWithAverage(gyroData, telemetryMetadata, 250)
	return gpsDataSamples, gyroDataSamples
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

		_, err = file.Read(buffer[bufferPos : bufferPos+chunkSize])
		if err != nil {
			fmt.Printf("Error reading at offset %d: %v\n", offset, err)
			return nil, err
		}
		bufferPos += chunkSize
	}
	return buffer, nil
}

func getUnixTimestamp(creationTimeV0 uint32) int64 {
	// MP4 Epoch starts at 1904-01-01, Unix Epoch starts at 1970-01-01
	mp4EpochOffset := int64(2082844800)

	// Convert to Unix timestamp
	return (int64(creationTimeV0) - mp4EpochOffset) * 1000
}

func assignTimestampsToGps(gpsData []GPS9, telemetryMetadata *TelemetryMetadata) []GPSSample {
	var gpsSamples []GPSSample
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for i := 0; i < int(timeToSample.SampleCount); i++ {

			if sampleIndex >= uint32(len(gpsData)) {
				break
			}

			sampleTime := telemetryMetadata.CreationTime + int64(sampleScaleTime*1000/telemetryMetadata.TimeScale)
			gpsSamples = append(gpsSamples, GPSSample{GPS9: gpsData[sampleIndex], TimeSample: TimeSample{TimeStamp: sampleTime}})
			sampleIndex++
			sampleScaleTime += timeToSample.SampleDelta
		}
	}

	return gpsSamples
}

// todo: refactor
func assignTimestampsToGyroWithAverage(
	gyroData [][]Gyroscope,
	telemetryMetadata *TelemetryMetadata,
	downsampleIntervalMs uint32,
) []GyroSample {
	var gyroSamples []GyroSample
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	var accumulatedGyro Gyroscope
	var accumulatedTime int64 = 0
	var count uint32 = 0
	var lastSampleScaleTime int64 = 0

	// Precompute factor to avoid repeated calculation
	downsampleScaleThreshold := int64(telemetryMetadata.TimeScale * downsampleIntervalMs / 1000)

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range int(timeToSample.SampleCount) {
			if sampleIndex >= uint32(len(gyroData)) {
				break
			}

			currentGyroSamples := gyroData[sampleIndex]
			sampleCount := uint32(len(currentGyroSamples))

			for _, gyro := range currentGyroSamples {
				accumulatedTime += int64(sampleScaleTime)

				// Accumulate gyro values
				accumulatedGyro.X += gyro.X
				accumulatedGyro.Y += gyro.Y
				accumulatedGyro.Z += gyro.Z
				count++

				// Check if enough time has passed to downsample
				if int64(sampleScaleTime)-lastSampleScaleTime >= downsampleScaleThreshold {
					avgGyro := averageGyro(accumulatedGyro, count)
					avgTime := calculateAverageTime(telemetryMetadata.CreationTime, accumulatedTime, count, telemetryMetadata.TimeScale)

					gyroSamples = append(gyroSamples, GyroSample{
						Gyroscope: avgGyro,
						TimeSample: TimeSample{
							TimeStamp: avgTime,
						},
					})

					// Reset accumulators
					accumulatedGyro = Gyroscope{}
					lastSampleScaleTime = int64(sampleScaleTime)
					accumulatedTime = 0
					count = 0
				}

				// Increment scaled time based on sample delta
				sampleScaleTime += timeToSample.SampleDelta / sampleCount
			}
			sampleIndex++
		}
	}

	return gyroSamples
}

// Helper: Compute average Gyroscope reading
func averageGyro(accumulated Gyroscope, count uint32) Gyroscope {
	return Gyroscope{
		X: accumulated.X / float32(count),
		Y: accumulated.Y / float32(count),
		Z: accumulated.Z / float32(count),
	}
}

// Helper: Compute average timestamp
func calculateAverageTime(creationTime int64, accumulatedTime int64, count uint32, timeScale uint32) int64 {
	averageScaleTime := accumulatedTime / int64(count)
	return creationTime + 1000*(averageScaleTime/int64(timeScale))
}
