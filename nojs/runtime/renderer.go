package runtime

import "github.com/ForgeLogic/nojs/vdom"

// Renderer defines the minimal set of runtime operations used by generated Render() code.
// This interface has NO build tags, making it available to both WASM and native test builds.
// This allows the AOT compiler to generate identical Render() signatures for both environments.
type Renderer interface {
	// RenderChild is used by generated code to render child components.
	// The key parameter uniquely identifies the component instance for state preservation.
	RenderChild(key string, childWithProps Component) *vdom.VNode

	// ReRender requests that the renderer re-run the render cycle.
	// Used by StateHasChanged() when component state changes.
	ReRender()

	// ReRenderSlot patches only the BodyContent slot of a layout,
	// preserving the layout instance and its state.
	// Called when a page component (inside a layout's slot) calls StateHasChanged().
	// It diffs the entire parent layout's VDOM; only changed content is patched.
	ReRenderSlot(slotParent Component) error

	// Navigate performs client-side navigation to the given path.
	// Used by Link components and programmatic navigation.
	Navigate(path string) error
}
