package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"gopro/gpmfParser"
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
			gpsData, gyroData, faceData := gpmfParser.ExtractTelemetryData(buf, false)

			// Resolve the Promise with the GPS data
			resolve.Invoke(map[string]interface{}{
				"gpsData":  convertGPSToJS(gpsData),
				"gyroData": convertGyroToJS(gyroData),
				"faceData": convertFaceToJS(faceData),
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
func convertGPSToJS(gpsData []gpmfParser.TimedGPS) js.Value {
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

func convertGyroToJS(gpsData []gpmfParser.TimedGyro) js.Value {
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

func convertFaceToJS(faceData []gpmfParser.TimedFace) js.Value {
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
