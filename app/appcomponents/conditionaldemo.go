//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
)

// ConditionalDemo demonstrates various conditional rendering scenarios
type ConditionalDemo struct {
	runtime.ComponentBase
	ShowWelcome bool
	IsLoggedIn  bool
	IsAdmin     bool
	HasError    bool
}

func (c *ConditionalDemo) ToggleWelcome() {
	c.ShowWelcome = !c.ShowWelcome
	c.StateHasChanged()
}

func (c *ConditionalDemo) ToggleLogin() {
	c.IsLoggedIn = !c.IsLoggedIn
	if !c.IsLoggedIn {
		c.IsAdmin = false // Can't be admin if not logged in
	}
	c.StateHasChanged()
}

func (c *ConditionalDemo) ToggleAdmin() {
	if c.IsLoggedIn {
		c.IsAdmin = !c.IsAdmin
		c.StateHasChanged()
	}
}

func (c *ConditionalDemo) ToggleError() {
	c.HasError = !c.HasError
	c.StateHasChanged()
}
