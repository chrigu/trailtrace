package parser

type Luma struct {
	Luminance uint8
}

func ParseLumaData(klvs []KLV) [][]Luma {
	return extractSensorData(klvs,
		"Average luminance",
		extractLumaData)
}

func extractLumaData(klv KLV) []Luma {

	var payload [][]uint8

	for _, child := range klv.Children {

		switch child.FourCC {
		case "YAVG":
			payload = readPayload(child).([][]uint8)

		default:
			continue
		}
	}

	lumminanceData := make([]Luma, len(payload))
	for i := range payload {
		lumminanceData[i] = Luma{
			Luminance: payload[i][0],
		}
	}

	return lumminanceData
}
