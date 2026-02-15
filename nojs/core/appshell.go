//go:build js || wasm
// +build js wasm

package core

import (
	"fmt"

	"github.com/vcrobe/nojs/console"
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/vdom"
)

// AppShell is a stable root component that holds persistent layouts (app shell)
// and swaps only the BodyContent slot when navigation occurs. This preserves
// layout instances and their internal state across navigations including sublayouts.
type AppShell struct {
	runtime.ComponentBase

	// persistent layout instance (app shell)
	persistentLayout runtime.Component

	// current chain of component instances (all from router, volatile)
	currentChain []runtime.Component
	currentKey   string
}

// NewAppShell creates a new AppShell with the given persistent layout component.
// The layout should implement the slot convention: a BodyContent []*vdom.VNode field.
func NewAppShell(persistentLayout runtime.Component) *AppShell {
	return &AppShell{
		persistentLayout: persistentLayout,
		currentChain:     make([]runtime.Component, 0),
	}
}

// SetPage replaces the volatile chain of component instances and triggers a re-render.
// The chain includes components from the router (from pivot onwards).
// When pivot > 0, the chain doesn't include the persistent layout (it's preserved).
// The RenderChild mechanism ensures layouts are reused efficiently,
// and VDOM patching only updates what changed.
func (a *AppShell) SetPage(chain []runtime.Component, key string) {
	console.Log("[AppShell.SetPage] Called with", len(chain), "components, key:", key)
	if len(chain) > 0 {
		console.Log("[AppShell.SetPage] First component type:", fmt.Sprintf("%T", chain[0]))
	}

	// If the chain doesn't include persistentLayout at index 0, prepend it
	// (this happens when pivot > 0 and layouts are preserved)
	if len(chain) == 0 || chain[0] != a.persistentLayout {
		console.Log("[AppShell.SetPage] Prepending persistentLayout to chain")
		fullChain := make([]runtime.Component, 0, len(chain)+1)
		fullChain = append(fullChain, a.persistentLayout)
		fullChain = append(fullChain, chain...)
		a.currentChain = fullChain
	} else {
		a.currentChain = chain
	}
	a.currentKey = key

	// Trigger a re-render of AppShell. RenderChild will reuse persistentLayout instance,
	// and VDOM patching will only update the changed slot content.
	console.Log("[AppShell.SetPage] Calling StateHasChanged")
	a.StateHasChanged()
}

// Render composes the persistent layout with the current component chain.
// The chain is linked here (child components injected into parent slots) because
// the router skips chain linking when using the AppShell pattern.
// We process the chain bottom-up, rendering each component and injecting it into
// its parent's BodyContent slot, then inject the root of the chain into the persistent layout.
func (a *AppShell) Render(r runtime.Renderer) *vdom.VNode {
	console.Log("[AppShell.Render] Called, chain length:", len(a.currentChain))

	// Ensure renderer is injected into children so StateHasChanged / Navigate work.
	type rendererSetter interface {
		SetRenderer(runtime.Renderer)
	}

	// Ensure persistent layout has renderer
	if a.persistentLayout != nil {
		if rs, ok := interface{}(a.persistentLayout).(rendererSetter); ok {
			rs.SetRenderer(r)
		}
	}

	// Link the chain: inject each child into parent's BodyContent slot
	// This is necessary because the router skips chain linking when using AppShell pattern
	var slotChildren []*vdom.VNode
	if len(a.currentChain) > 0 {
		console.Log("[AppShell.Render] Processing chain")
		// Skip persistent layout if it's at index 0 (AppShell manages it separately)
		chainIndex := 0
		if a.currentChain[0] == a.persistentLayout {
			console.Log("[AppShell.Render] First component is persistentLayout, skipping")
			chainIndex = 1
		}

		// Link the chain from parent to child (bottom-up: leaf → root)
		// Start from the last component and work backwards
		for i := len(a.currentChain) - 1; i > chainIndex; i-- {
			child := a.currentChain[i]
			parent := a.currentChain[i-1]

			// Give renderer to child component if possible
			if rs, ok := interface{}(child).(rendererSetter); ok {
				rs.SetRenderer(r)
			}

			// Render child to VDOM and inject into parent's slot
			slotKey := fmt.Sprintf("slot-chain-%d-%T-%p", i, child, child)
			childVNode := r.RenderChild(slotKey, child)
			if childVNode != nil {
				console.Log("[AppShell.Render] Linking", fmt.Sprintf("%T", child), "into", fmt.Sprintf("%T", parent))
				// Use duck typing to set slot content - any layout with SetBodyContent method
				if layout, ok := parent.(interface{ SetBodyContent([]*vdom.VNode) }); ok {
					layout.SetBodyContent([]*vdom.VNode{childVNode})
				}
			}
		}

		// Now render the first non-layout component in the chain
		// (which now has all its children properly linked)
		if chainIndex < len(a.currentChain) {
			rootComponent := a.currentChain[chainIndex]
			console.Log("[AppShell.Render] Rendering root component at index", chainIndex, "type:", fmt.Sprintf("%T", rootComponent))

			// Give renderer to root component if possible
			if rs, ok := interface{}(rootComponent).(rendererSetter); ok {
				rs.SetRenderer(r)
			}

			// Use RenderChild to track sublayouts and pages for efficient caching/patching.
			// The key includes the type and pointer so different component types don't collide,
			// but the same preserved instance (e.g., AdminLayout across /admin → /admin/settings)
			// gets reused and efficiently patched.
			slotKey := fmt.Sprintf("slot-root-%T-%p", rootComponent, rootComponent)
			childVNode := r.RenderChild(slotKey, rootComponent)
			if childVNode != nil {
				slotChildren = []*vdom.VNode{childVNode}
			}
		}
	}

	// Inject into layout's BodyContent slot (compiler-generated field)
	// Layouts follow the single-slot convention: BodyContent []*vdom.VNode
	if a.persistentLayout != nil {
		// assign slot directly; generated ApplyProps will preserve state on instance reuse
		if layout, ok := a.persistentLayout.(interface{ SetBodyContent([]*vdom.VNode) }); ok {
			layout.SetBodyContent(slotChildren)
		}
		// Use RenderChild to render persistentLayout so the renderer caches it for efficient patching
		// on subsequent navigations. The key "persistent-layout" identifies this component.
		return r.RenderChild("persistent-layout", a.persistentLayout)
	}

	// Fallback: if no layout, render the first non-layout component from chain
	if len(a.currentChain) > 0 {
		// Skip persistent layout if it's at index 0
		chainIndex := 0
		if a.currentChain[0] == a.persistentLayout {
			chainIndex = 1
		}
		if chainIndex < len(a.currentChain) {
			rootComponent := a.currentChain[chainIndex]
			slotKey := fmt.Sprintf("slot-root-%T-%p", rootComponent, rootComponent)
			return r.RenderChild(slotKey, rootComponent)
		}
	}

	// Empty fallback
	return vdom.NewVNode("div", nil, nil, "")
}
