//go:build js || wasm
// +build js wasm

package runtime

// ComponentBase is a struct that components can embed to gain access to the
// StateHasChanged method, which triggers a UI re-render.
type ComponentBase struct {
	renderer *Renderer
}

// SetRenderer is called by the framework's runtime to inject a reference
// to the renderer, enabling StateHasChanged. This method should not be
// called by user code.
func (b *ComponentBase) SetRenderer(r *Renderer) {
	b.renderer = r
}

// StateHasChanged signals to the framework that the component's state has
// been updated and the UI should be re-rendered to reflect the changes.
func (b *ComponentBase) StateHasChanged() {
	if b.renderer == nil {
		println("StateHasChanged called, but renderer is nil (component not mounted?)")
		return
	}
	// Trigger a re-render of the root component.
	b.renderer.ReRender()
}
