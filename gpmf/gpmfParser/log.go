//go:build !js
// +build !js

package gpmfParser

import "fmt"

// Log prints messages in CLI builds but is disabled in WASM.
func log(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}
