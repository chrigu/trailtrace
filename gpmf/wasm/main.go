package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"gopro/telemetry"
)

func main() {
	js.Global().Set("processFile", js.FuncOf(processFile))
	select {}
}

func processFile(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		fmt.Println("Error: No file specified")
		return js.Null()
	}

	file := args[0]

	// Create a new Promise
	promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, resolveArgs []js.Value) any {
		resolve := resolveArgs[0]

		// Use FileReader to read the content
		fileReader := js.Global().Get("FileReader").New()
		fileReader.Set("onload", js.FuncOf(func(this js.Value, p []js.Value) any {
			data := p[0].Get("target").Get("result")
			buffer := js.Global().Get("Uint8Array").New(data)

			// Convert Uint8Array to Go byte slice
			byteSlice := make([]byte, buffer.Length())
			js.CopyBytesToGo(byteSlice, buffer)

			fmt.Printf("Buffer Length: %d bytes\n", len(byteSlice))

			// Create a bytes.Reader from the byte slice
			buf := bytes.NewReader(byteSlice)

			// Extract GPS data
			gpsData, gyroData, faceData, lumaData, colorData := telemetry.ExtractTelemetryData(buf, false)

			// Resolve the Promise with the GPS data
			resolve.Invoke(map[string]interface{}{
				"gpsData":  convertGPSToJS(gpsData),
				"gyroData": convertGyroToJS(gyroData),
				"faceData": convertFaceToJS(faceData),
				"lumaData": convertLumaToJS(lumaData),
				"hueData":  convertColorsToJS(colorData),
			})
			return nil
		}))

		// Read file as ArrayBuffer
		fileReader.Call("readAsArrayBuffer", file)

		return nil // Promise will resolve later
	}))

	return promise
}

// Converts []GPS9 to a JavaScript array of objects
func convertGPSToJS(gpsData []telemetry.TimedGPS) js.Value {
	jsArray := js.Global().Get("Array").New(len(gpsData))
	for i, coord := range gpsData {
		// Create a JavaScript object for each GPS9 struct
		jsCoord := js.Global().Get("Object").New()
		jsCoord.Set("latitude", coord.Latitude)
		jsCoord.Set("longitude", coord.Longitude)
		jsCoord.Set("altitude", coord.Altitude)
		jsCoord.Set("timestamp", coord.TimeStamp)
		jsArray.SetIndex(i, jsCoord)
	}
	return jsArray
}

func convertGyroToJS(gpsData []telemetry.TimedGyro) js.Value {
	jsArray := js.Global().Get("Array").New(len(gpsData))
	for i, coord := range gpsData {
		// Create a JavaScript object for each GPS9 struct
		jsCoord := js.Global().Get("Object").New()
		jsCoord.Set("x", coord.X)
		jsCoord.Set("y", coord.Y)
		jsCoord.Set("z", coord.Z)
		jsCoord.Set("timestamp", coord.TimeStamp)
		jsArray.SetIndex(i, jsCoord)
	}
	return jsArray
}

func convertFaceToJS(faceData []telemetry.TimedFace) js.Value {
	jsArray := js.Global().Get("Array").New(len(faceData))
	for i, face := range faceData {
		jsFace := js.Global().Get("Object").New()
		jsFace.Set("confidence", face.Confidence)
		jsFace.Set("id", face.ID)
		jsFace.Set("x", face.X)
		jsFace.Set("y", face.Y)
		jsFace.Set("w", face.W)
		jsFace.Set("h", face.H)
		jsFace.Set("smile", face.Smile)
		jsFace.Set("blink", face.Blink)
		jsFace.Set("timestamp", face.TimeStamp)
		jsArray.SetIndex(i, jsFace)
	}
	return jsArray
}

func convertLumaToJS(lumaData []telemetry.TimedLuma) js.Value {
	jsArray := js.Global().Get("Array").New(len(lumaData))
	for i, luma := range lumaData {
		jsLuma := js.Global().Get("Object").New()
		jsLuma.Set("luma", luma.Luminance)
		jsLuma.Set("timestamp", luma.TimeStamp)
		jsArray.SetIndex(i, jsLuma)
	}
	return jsArray
}

func convertColorsToJS(colorData []telemetry.TimedColor) js.Value {
	jsArray := js.Global().Get("Array").New(len(colorData))
	for i, color := range colorData {
		jsColor := js.Global().Get("Object").New()
		jsRed := js.Global().Get("Object").New()
		jsRed.Set("hue", color.Red.Hue)
		jsRed.Set("weight", color.Red.Weight)

		jsGreen := js.Global().Get("Object").New()
		jsGreen.Set("hue", color.Green.Hue)
		jsGreen.Set("weight", color.Green.Weight)

		jsBlue := js.Global().Get("Object").New()
		jsBlue.Set("hue", color.Blue.Hue)
		jsBlue.Set("weight", color.Blue.Weight)

		jsColor.Set("red", jsRed)
		jsColor.Set("green", jsGreen)
		jsColor.Set("blue", jsBlue)
		jsColor.Set("timestamp", color.TimeStamp)

		jsArray.SetIndex(i, jsColor)
	}
	return jsArray
}
