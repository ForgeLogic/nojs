//go:build js || wasm
// +build js wasm

package runtime

import (
	"github.com/vcrobe/nojs/vdom"
)

// Renderer is the core of the framework's runtime. It manages the component instance tree.
type Renderer struct {
	instances map[string]Component
	root      Component
	mountID   string
}

// NewRenderer creates a new runtime renderer.
func NewRenderer(root Component, mountID string) *Renderer {
	return &Renderer{
		instances: make(map[string]Component),
		root:      root,
		mountID:   mountID,
	}
}

// RenderRoot starts the rendering process for the entire application.
func (r *Renderer) RenderRoot() {
	// On each root render, we build the VDOM tree from the root component.
	vdomTree := r.root.Render(r)

	// Simple strategy: clear the mount point and render fresh.
	vdom.Clear(r.mountID)
	vdom.RenderToSelector(r.mountID, vdomTree)
}

// RenderChild is called by compiler-generated code to render a child component.
// It handles the core logic of instance creation and reuse.
func (r *Renderer) RenderChild(key string, childWithProps Component) *vdom.VNode {
	instance, ok := r.instances[key]
	if !ok {
		// First time seeing this component at this location, so store the new instance.
		instance = childWithProps
		r.instances[key] = instance
	} else {
		// We have seen this component before. We need to apply the new props
		// to the *existing* instance. For now, since the compiler creates a new
		// struct literal each time, we will overwrite the instance.
		// A more advanced implementation would use reflection to update fields
		// on the existing `instance` without replacing it. For now, we will
		// just replace it to keep the logic simple, but this means state is not preserved.
		// To truly preserve state, we'd need to update the existing instance's fields.
		// For now, let's just re-assign.
		r.instances[key] = childWithProps
		instance = childWithProps
	}

	// Now, render the child (either the new or reused one).
	return instance.Render(r)
}
