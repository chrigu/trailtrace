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
	fmt.Println("Processing file...", this, args, len(args))
	if len(args) < 1 {
		fmt.Println("Error: No file specified")
		return nil
	}

	// Get file object properties
	file := args[0]
	fileName := file.Get("name").String()
	fileSize := file.Get("size").Int()
	fileType := file.Get("type").String()

	// Log debugging information about the file
	fmt.Printf("File Name: %s\n", fileName)
	fmt.Printf("File Size: %d bytes\n", fileSize)
	fmt.Printf("File Type: %s\n", fileType)

	// Use FileReader to read the content
	fileReader := js.Global().Get("FileReader").New()
	fileReader.Set("onload", js.FuncOf(func(this js.Value, p []js.Value) any {
		data := p[0].Get("target").Get("result")          // Fix index error
		buffer := js.Global().Get("Uint8Array").New(data) // Fix Uint9Array typo

		// Convert Uint8Array to Go byte slice
		byteSlice := make([]byte, buffer.Length())
		js.CopyBytesToGo(byteSlice, buffer)

		fmt.Printf("Buffer Length: %d bytes\n", len(byteSlice))

		// Prevent out-of-bounds slice
		sliceLength := len(byteSlice)
		if sliceLength > 100 {
			sliceLength = 100
		}
		fmt.Printf("First %d bytes: %x\n", sliceLength, byteSlice[:sliceLength])

		// Create a bytes.Reader from the byte slice
		buf := bytes.NewReader(byteSlice)

		// Call telemetry data extraction function
		gpsData := gpmfParser.ExtractTelemetryData(buf)
		return nil
	}))

	// Use FileReader's readAsArrayBuffer to get binary content
	fileReader.Call("readAsArrayBuffer", file)
	return nil
}
