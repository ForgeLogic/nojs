//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
)

type SimpleMessage struct {
	runtime.ComponentBase
	Counter   int
	FirstProp string
}

func (r *SimpleMessage) Increment() {
	r.Counter++

	println("Called Increment: Counter is now", r.Counter)
	println("FirstProp value is %s", r.FirstProp)

	r.StateHasChanged()
}
