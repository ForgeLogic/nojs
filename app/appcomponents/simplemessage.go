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
	IsSaved   bool
	IsLoading bool
}

func (r *SimpleMessage) Increment() {
	r.Counter++

	println("Called Increment: Counter is now", r.Counter)
	println("FirstProp value is", r.FirstProp)

	r.IsLoading = true

	r.StateHasChanged()
}
