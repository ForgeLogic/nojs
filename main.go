//go:build js || wasm
// +build js wasm

package main

import (
	"github.com/vcrobe/nojs/appcomponents" // Assuming components are in this package
	"github.com/vcrobe/nojs/runtime"
)

func main() {
	// The user creates the root component instance.
	// The generated .gt.go files for the 'appcomponents' package
	// must be compiled as part of the main application.
	app := appcomponents.App{}

	// Create the runtime renderer, passing it the root component.
	renderer := runtime.NewRenderer(&app, "#app")

	// Tell the renderer to perform the initial render.
	renderer.RenderRoot()

	// Keep the Go program running.
	select {}
}
