package mp4

import (
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
