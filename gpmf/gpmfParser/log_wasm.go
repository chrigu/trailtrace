//go:build js && wasm
// +build js,wasm

package gpmfParser

// Log is a no-op in WASM builds
func log(v ...interface{}) {
	// Do nothing (silenced logs)
}
