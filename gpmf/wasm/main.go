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
			gpsData, gyroData, faceData, lumaData, colorData, sceneData := telemetry.ExtractTelemetryData(buf, false)

			// Resolve the Promise with the GPS data
			resolve.Invoke(map[string]interface{}{
				"gpsData":   convertGPSToJS(gpsData),
				"gyroData":  convertGyroToJS(gyroData),
				"faceData":  convertFaceToJS(faceData),
				"lumaData":  convertLumaToJS(lumaData),
				"hueData":   convertHuesToJS(colorData),
				"sceneData": convertSceneToJS(sceneData),
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

func convertHuesToJS(hueData []telemetry.TimedHue) js.Value {
	jsArray := js.Global().Get("Array").New(len(hueData))
	for i, timedHue := range hueData {
		// Create a JavaScript object for each TimedHue
		jsHue := js.Global().Get("Object").New()

		// Create an array for the hues
		jsHues := js.Global().Get("Array").New(len(timedHue.Hues))

		// Process each hue in the Hues array
		for j, hue := range timedHue.Hues {
			// Create a JavaScript object for each hue
			jsHueObj := js.Global().Get("Object").New()
			jsHueObj.Set("hue", hue.Hue)
			jsHueObj.Set("weight", hue.Weight)

			// Add to the hues array
			jsHues.SetIndex(j, jsHueObj)
		}

		// Set the hues array and timestamp
		jsHue.Set("hues", jsHues)
		jsHue.Set("timestamp", timedHue.TimeStamp)

		// Add to the main array
		jsArray.SetIndex(i, jsHue)
	}
	return jsArray
}

func convertSceneToJS(sceneData []telemetry.TimedScene) js.Value {
	jsArray := js.Global().Get("Array").New(len(sceneData))
	for i, scene := range sceneData {
		jsScene := js.Global().Get("Object").New()
		jsScene.Set("scene", string([]byte(scene.Scene.Type)))
		jsScene.Set("probability", scene.Scene.Prob)
		jsScene.Set("timestamp", scene.TimeStamp)
		jsArray.SetIndex(i, jsScene)
	}
	return jsArray
}
