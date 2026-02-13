//go:build js || wasm

package pages

import (
	"github.com/vcrobe/app/internal/app/components/shared/modal"
	"github.com/vcrobe/app/internal/app/context"
	"github.com/vcrobe/nojs/runtime"
)

// HomePage is the component rendered for the "/" route.
type HomePage struct {
	runtime.ComponentBase

	MainLayoutCtx *context.MainLayoutCtx

	Years   []int
	Counter int

	// The parent *must* control the visibility state.
	IsMyModalVisible bool
	// We can store a message to show the result
	LastModalResult string
}

func (h *HomePage) OnInit() {
	h.Years = []int{120, 300}
}

func (c *HomePage) UpdateTitle() {
	c.MainLayoutCtx.SetTitle("Updated title from HomePage")
	// Note: SettingsPage doesn't need to call StateHasChanged()
	// for the header; the SetTitle method handles it for the Layout!
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

func (h *HomePage) IncrementCounter() {
	h.Counter++
	h.StateHasChanged()
}
