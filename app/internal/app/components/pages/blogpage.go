//go:build js || wasm
// +build js wasm

package pages

import (
	"github.com/vcrobe/nojs/events"
	"github.com/vcrobe/nojs/runtime"
)

// BlogPage is the component rendered for the "/blog/{year}" route.

type BlogPage struct {
	runtime.ComponentBase

	Year int
}

// NavigateToHome demonstrates navigation back to the home page using the `<a>` tag.
func (a *BlogPage) NavigateToHome(e events.ClickEventArgs) {
	e.PreventDefault()
	if err := a.Navigate("/"); err != nil {
		println("Navigation error:", err.Error())
	}
}
