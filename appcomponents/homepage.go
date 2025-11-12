//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/events"
	"github.com/vcrobe/nojs/runtime"
)

// HomePage is the component rendered for the "/" route.
type HomePage struct {
	runtime.ComponentBase
	Years []int
}

func (h *HomePage) OnInit() {
	h.Years = []int{2025, 1930, 2000, 2010}
}

// NavigateToAbout handles navigation to the about page
func (h *HomePage) NavigateToAbout(e events.ClickEventArgs) {
	e.PreventDefault()
	if err := h.Navigate("/about"); err != nil {
		println("Navigation error:", err.Error())
	}
}

func (h *HomePage) NavigateToBlogs(e events.ClickEventArgs) {
	e.PreventDefault()
	if err := h.Navigate("/blog/2025"); err != nil {
		println("Navigation error:", err.Error())
	}
}
