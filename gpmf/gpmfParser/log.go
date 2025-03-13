//go:build !js
// +build !js

package gpmfParser

import "fmt"

// Log prints messages in CLI builds but is disabled in WASM.
func log(v ...interface{}) {
	fmt.Print(v...)
}
