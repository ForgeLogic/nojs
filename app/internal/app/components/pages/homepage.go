//go:build js || wasm
// +build js wasm

package pages

import (
	"time"

	"github.com/vcrobe/app/internal/app/components/shared/modal"
	"github.com/vcrobe/nojs/runtime"
)

// HomePage is the component rendered for the "/" route.
type HomePage struct {
	runtime.ComponentBase

	disposed bool

	Years []int

	// The parent *must* control the visibility state.
	IsMyModalVisible bool
	// We can store a message to show the result
	LastModalResult string
	IsLoading       bool
	LoadedData      string
}

func (h *HomePage) OnInit() {
	h.Years = []int{120, 300}

	// // Start the async loading operation
	h.IsLoading = true
	h.LoadedData = "Loading..."

	// // Run async operation in a goroutine
	// go func() {
	// 	// Simulate network delay
	// 	time.Sleep(2 * time.Second)

	// 	// Update component state after async operation completes
	// 	h.IsLoading = false
	// 	h.LoadedData = "Data loaded successfully!"

	// 	println("Async operation complete: ", h.LoadedData)

	// 	// CRITICAL: Call StateHasChanged() to trigger re-render after async state update
	// 	h.StateHasChanged()
	// }()
}

func (h *HomePage) OnDispose() {
	// Clean up any resources, cancel timers, etc.
	// This is called when the component is removed from the UI.
	h.disposed = true
}

func (h *HomePage) OnParametersSet() {
	// OnParametersSet is called before every render, including the first one.
	// For simple cases where you only need one-time initialization,
	// use OnInit() instead.

	// Kick off async work whenever parameters change.
	if !h.IsLoading {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if h.disposed {
					return
				}
				h.IsLoading = false
				h.LoadedData = "Something went wrong"
				h.StateHasChanged()
			}
		}()

		time.Sleep(2 * time.Second) // simulate work

		if h.disposed {
			return // component is gone; skip state update
		}

		h.IsLoading = false
		h.LoadedData = "Data loaded successfully!"

		println("Async operation complete: ", h.LoadedData)

		h.StateHasChanged()
	}()
}

// ShowTheModal is called by our button.
func (c *HomePage) ShowTheModal() {
	c.IsMyModalVisible = true
	c.LastModalResult = "Modal is open..."
	// CRITICAL (Rule 6): We changed state, so we *must* call StateHasChanged().
	c.StateHasChanged()
}

// HandleModalClose matches the 'OnClose' prop signature.
// This is how the dialog communicates back to the parent.
func (c *HomePage) HandleModalClose(result modal.ModalResult) {
	c.IsMyModalVisible = false // Hide the dialog

	if result == modal.Ok {
		c.LastModalResult = "You clicked OK!"
	} else {
		c.LastModalResult = "You clicked Cancel."
	}

	c.StateHasChanged()
}
