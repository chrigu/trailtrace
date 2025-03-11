package main

import (
	"fmt"
	"os"
	"bytes"
	"github.com/abema/go-mp4"
)
// Extracts raw box data and writes it to a file
func extractBoxData(h *mp4.ReadHandle, fileName string) error {
	buf := new(bytes.Buffer)

	// Read raw box data into the buffer
	n, err := h.ReadData(buf)
	if err != nil {
		return fmt.Errorf("error reading data from %s: %w", fileName, err)
	}
  fmt.Printf("Bytes read: %d\n", n)

	// Write to a file
	return writeBinaryToFile(fileName, buf.Bytes())
}

// Function to write binary data to a file
func writeBinaryToFile(fileName string, data []byte) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}
	defer outFile.Close()

	_, err = outFile.Write(data)
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}

	fmt.Println("Saved data to:", fileName)
	return nil
}