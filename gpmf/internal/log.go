//go:build !js
// +build !js

package internal

import "fmt"

// Log prints messages in CLI builds but is disabled in WASM.
func Log(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}
