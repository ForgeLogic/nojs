//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/events"
	"github.com/vcrobe/nojs/runtime"
)

// AboutPage is the component rendered for the "/about" route.
type AboutPage struct {
	runtime.ComponentBase
}

// NavigateToHome handles navigation to the home page
func (a *AboutPage) NavigateToHome(e events.ClickEventArgs) {
	println("NavigateToHome called!")
	e.PreventDefault()
	println("PreventDefault called, about to navigate...")
	if err := a.Navigate("/"); err != nil {
		println("Navigation error:", err.Error())
	} else {
		println("Navigation successful!")
	}
}
