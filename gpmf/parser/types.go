package parser

type KLV struct {
	FourCC   string
	DataType int
	DataSize uint32
	Repeat   uint32
	Payload  []byte
	Children []KLV
}
type GPS9 struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
}

// Rename or add Accelerometer struct
type Gyroscope struct {
	X float32
	Y float32
	Z float32
}

type Face struct {
	Confidence float32
	ID         int
	X          float32
	Y          float32
	W          float32
	H          float32
	Smile      float32
	Blink      float32
}
