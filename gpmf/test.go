package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func Test() {
	file, err := os.Open("../GX010025.MP4")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Define the binary pattern to search for
	pattern := []byte("BlackSTRM") // Example: Looking for the "moov" box
	bufferSize := 1024 // Read in chunks

	reader := bufio.NewReader(file)
	var offset int64 = 0

	for {
		// Read a chunk from the file
		buffer := make([]byte, bufferSize)
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}

		// Search for the pattern in the current buffer
		if idx := bytes.Index(buffer, pattern); idx != -1 {
			fmt.Printf("Found pattern at offset: %d\n", offset+int64(idx))
			break
		}

		// Update offset for next chunk
		offset += int64(n)
	}

	fmt.Println("Finished searching")
}
