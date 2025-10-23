//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/vdom"
)

// Card is a reusable layout component with content projection.
type Card struct {
	runtime.ComponentBase
	Title       string
	BodyContent []*vdom.VNode // Content slot - any exported field of type []*vdom.VNode becomes a slot
}
