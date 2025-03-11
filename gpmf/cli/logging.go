package main

import "fmt"

// Log prints messages in CLI builds but is disabled in WASM.
func Log(v ...interface{}) {
	fmt.Println(v...)
}
