//go:build js && wasm
// +build js,wasm

package internal

// Log is a no-op in WASM builds
func Log(v ...interface{}) {
	// Do nothing (silenced logs)
	// fmt.Print(v...)
}
