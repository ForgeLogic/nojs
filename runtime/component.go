//go:build js || wasm
// +build js wasm

package runtime

import (
	"github.com/vcrobe/nojs/vdom"
)

// Component interface defines the structure for all components in the framework.
// The Render method now accepts the runtime renderer to manage child components.
type Component interface {
	Render(r *Renderer) *vdom.VNode
	SetRenderer(r *Renderer)
}
