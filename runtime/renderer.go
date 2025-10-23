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
	prevVDOM  *vdom.VNode // Previous VDOM tree for patching
}

// NewRenderer creates a new runtime renderer.
func NewRenderer(root Component, mountID string) *Renderer {
	return &Renderer{
		instances: make(map[string]Component),
		root:      root,
		mountID:   mountID,
		prevVDOM:  nil,
	}
}

// RenderRoot starts the rendering process for the entire application.
func (r *Renderer) RenderRoot() {
	// On each root render, we build the VDOM tree from the root component.
	// Ensure the root has a reference to the renderer for StateHasChanged.
	if r.root != nil {
		r.root.SetRenderer(r)
	}
	newVDOM := r.root.Render(r)

	if r.prevVDOM == nil {
		// Initial render: clear and render fresh
		vdom.Clear(r.mountID)
		vdom.RenderToSelector(r.mountID, newVDOM)
	} else {
		// Subsequent renders: patch the existing DOM
		vdom.Patch(r.mountID, r.prevVDOM, newVDOM)
	}

	// Store the new VDOM tree for the next render cycle
	r.prevVDOM = newVDOM
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
		// We have seen this component before. Preserve the existing instance to keep state.
		// In the future, apply new props onto the existing instance instead of replacing it.
	}

	// Now, render the child (either the new or reused one).
	// Ensure the instance knows about the renderer so it can call StateHasChanged.
	instance.SetRenderer(r)
	return instance.Render(r)
}

// ReRender patches the DOM with minimal changes.
func (r *Renderer) ReRender() {
	r.RenderRoot()
}
