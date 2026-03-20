//go:build js || wasm

package components

import (
	"github.com/ForgeLogic/nojs/events"
	"github.com/ForgeLogic/nojs/runtime"
)

// HelloApp is a minimal nojs component.
// Fields on the struct are the component's state.
type HelloApp struct {
	runtime.ComponentBase

	Name    string
	HasName bool
}

// HandleInput is called on every keystroke in the text input.
// It updates the state and tells nojs to re-render.
func (c *HelloApp) HandleInput(e events.ChangeEventArgs) {
	c.Name = e.Value
	c.HasName = c.Name != ""
	c.StateHasChanged()
}
