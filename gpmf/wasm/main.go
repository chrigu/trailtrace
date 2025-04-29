package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"gopro/mp4"
	"gopro/telemetry"
)

func main() {
	js.Global().Set("processFile", js.FuncOf(processFile))
	js.Global().Set("exportGPMF", js.FuncOf(exportGPMF))
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
	for i, timedScene := range sceneData {
		jsScene := js.Global().Get("Object").New()
		// Create an array for all scenes
		jsScenes := js.Global().Get("Array").New(len(timedScene.Scenes))
		for j, scene := range timedScene.Scenes {
			jsSceneObj := js.Global().Get("Object").New()
			// Convert the FourCCScene to bytes and create a Uint8Array in JS
			jsSceneObj.Set("type", string(scene.Type))
			jsSceneObj.Set("probability", scene.Prob)
			jsScenes.SetIndex(j, jsSceneObj)
		}
		jsScene.Set("scenes", jsScenes)
		jsScene.Set("timestamp", timedScene.TimeStamp)
		jsArray.SetIndex(i, jsScene)
	}
	return jsArray
}

func exportGPMF(this js.Value, args []js.Value) any {
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

			gpmfRaw, _ := mp4.ExtractTelemetryFromMp4(buf)
			// if err != nil {
			// 	fmt.Printf("Error extracting GPMF: %v\n", err)
			// 	resolve.Invoke(js.Null())
			// 	return nil
			// }

			// Create a new Uint8Array in JavaScript
			jsGPMF := js.Global().Get("Uint8Array").New(len(gpmfRaw))
			// Copy the GPMF data to the JavaScript Uint8Array
			js.CopyBytesToJS(jsGPMF, gpmfRaw)

			// Resolve the Promise with the GPMF data
			resolve.Invoke(jsGPMF)
			return nil
		}))

		// Read file as ArrayBuffer
		fileReader.Call("readAsArrayBuffer", file)

		return nil
	}))

	return promise
}
