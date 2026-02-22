//go:build js || wasm
// +build js wasm

package runtime

import (
	"fmt"
	"sync"

	"github.com/ForgeLogic/nojs/vdom"
)

// Compile-time assertion to ensure the concrete RendererImpl implements the Renderer interface.
var _ Renderer = (*RendererImpl)(nil)

// RendererImpl is the concrete implementation of the Renderer interface.
// It manages the component instance tree and handles rendering lifecycle.
// All public methods are thread-safe and can be called from multiple goroutines.
type RendererImpl struct {
	mu                sync.Mutex // Protects all renderer state from concurrent access
	instances         map[string]Component
	initialized       map[string]bool   // Track which components have been initialized
	activeKeys        map[string]bool   // Track which components are active in the current render
	currentComponent  Component         // The currently active root component (set by router or directly)
	currentKey        string            // Key for component-level reconciliation (e.g., current route path)
	navManager        NavigationManager // Optional: router for client-side navigation
	mountID           string
	prevVDOM          *vdom.VNode               // Previous VDOM tree for patching
	instanceVDOMCache map[Component]*vdom.VNode // Track VDOM per component instance (for scoped updates)
	renderingStack    []Component               // Stack of components currently rendering (for scoped cache keys)
}

// NewRenderer creates a new runtime renderer.
// If navManager is provided, the renderer will support client-side routing.
// If navManager is nil, the renderer works without routing (useful for non-SPA apps).
func NewRenderer(navManager NavigationManager, mountID string) *RendererImpl {
	return &RendererImpl{
		instances:         make(map[string]Component),
		initialized:       make(map[string]bool),
		activeKeys:        make(map[string]bool),
		instanceVDOMCache: make(map[Component]*vdom.VNode),
		navManager:        navManager,
		mountID:           mountID,
		prevVDOM:          nil,
		renderingStack:    make([]Component, 0),
	}
}

// GetCurrentComponent returns the current root component being rendered.
// This is used by the router engine to access methods on the root component (e.g., AppShell).
func (r *RendererImpl) GetCurrentComponent() Component {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.currentComponent
}

// SetCurrentComponent sets the component to be rendered with an optional key.
// The key is used for component-level reconciliation (e.g., for router navigation).
// When the key changes, the entire component tree is replaced instead of patched.
// This is typically called by the router's onChange callback when navigation occurs.
// For non-routed apps, it can be called directly with a static component and empty key.
// This method is thread-safe.
func (r *RendererImpl) SetCurrentComponent(comp Component, key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentComponent = comp
	r.currentKey = key
}

// RenderRoot starts the rendering process for the entire application.
// This method is thread-safe and protected by a mutex.
func (r *RendererImpl) RenderRoot() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Reset activeKeys for this render cycle
	r.activeKeys = make(map[string]bool)

	// On each root render, we build the VDOM tree from the current component.
	// Ensure the component has a reference to the renderer for StateHasChanged and Navigate.
	if r.currentComponent != nil {
		r.currentComponent.SetRenderer(r)

		// Push root component onto rendering stack
		r.renderingStack = append(r.renderingStack, r.currentComponent)

		// Handle root component lifecycle
		if _, initialized := r.initialized["__root__"]; !initialized {
			// Call OnMount only once, before first render
			if mountable, ok := r.currentComponent.(Mountable); ok {
				r.callOnMount(mountable, "__root__")
			}
			r.initialized["__root__"] = true
		}

		// Call OnParametersSet before every render (including first)
		if paramReceiver, ok := r.currentComponent.(ParameterReceiver); ok {
			r.callOnParametersSet(paramReceiver, "__root__")
		}
	}

	newVDOM := r.currentComponent.Render(r)

	// Pop root component from rendering stack
	if len(r.renderingStack) > 0 {
		r.renderingStack = r.renderingStack[:len(r.renderingStack)-1]
	}

	// Attach the component key to the root VNode for reconciliation
	newVDOM.ComponentKey = r.currentKey

	if r.prevVDOM == nil {
		// Initial render: clear and render fresh
		vdom.Clear(r.mountID, nil)
		vdom.RenderToSelector(r.mountID, newVDOM)
	} else {
		// Check if component key changed (e.g., router navigation)
		if r.prevVDOM.ComponentKey != newVDOM.ComponentKey {
			// Component key changed - replace entire tree
			vdom.Clear(r.mountID, r.prevVDOM)
			vdom.RenderToSelector(r.mountID, newVDOM)

			// Call OnUnmount on old root component
			if unmountable, ok := r.currentComponent.(Unmountable); ok {
				r.callOnUnmount(unmountable, "__root__")
			}

			// Reset initialization tracking for fresh component lifecycle
			r.initialized = make(map[string]bool)
		} else {
			// Same key - patch normally
			vdom.Patch(r.mountID, r.prevVDOM, newVDOM)
		}
	}

	// Store the new VDOM tree for the next render cycle
	r.prevVDOM = newVDOM

	// Cache VDOM for the current component (for scoped slot updates)
	r.instanceVDOMCache[r.currentComponent] = newVDOM

	// Clean up components that were not rendered in this cycle
	r.cleanupUnmountedComponents()
}

// RenderChild is called by compiler-generated code to render a child component.
// It handles the core logic of instance creation and reuse.
// Uses a composite key that includes the parent component context to avoid collisions
// when multiple parent components render children with the same logical key.
func (r *RendererImpl) RenderChild(key string, childWithProps Component) *vdom.VNode {
	// Create a globally unique key by including the parent component's pointer
	// This prevents collisions when multiple parents render children with the same key
	globalKey := key
	if len(r.renderingStack) > 0 {
		parent := r.renderingStack[len(r.renderingStack)-1]
		globalKey = fmt.Sprintf("%p:%s", parent, key)
	}

	// Mark this component as active in the current render cycle
	r.activeKeys[globalKey] = true

	instance, exists := r.instances[globalKey]
	isFirstRender := false

	if !exists {
		// First time seeing this component at this location, so store the new instance.
		instance = childWithProps
		r.instances[globalKey] = instance
		isFirstRender = true
	} else {
		// We have seen this component before. Preserve the existing instance to keep state.
		// Apply new props from childWithProps to the existing instance.
		if updater, ok := instance.(PropUpdater); ok {
			println("[RenderChild] Found cached component, calling ApplyProps for key:", globalKey)
			updater.ApplyProps(childWithProps)
		}
	}

	// Now, render the child (either the new or reused one).
	// Ensure the instance knows about the renderer so it can call StateHasChanged.
	instance.SetRenderer(r)

	// Track slot parent relationship (child is inside parent's []*vdom.VNode slot)
	// This enables scoped re-renders when child calls StateHasChanged()
	// We track this by checking if the current component is a layout (has a parent field)
	if parentComponent, ok := childWithProps.(interface{ SetSlotParent(Component) }); ok {
		if r.currentComponent != nil {
			parentComponent.SetSlotParent(r.currentComponent)
		}
	}

	// Call lifecycle methods in the correct order
	if isFirstRender {
		// Call OnMount only once, before first render
		if mountable, ok := instance.(Mountable); ok {
			r.callOnMount(mountable, globalKey)
		}
		r.initialized[globalKey] = true
	}

	// Call OnParametersSet before every render (including first)
	if paramReceiver, ok := instance.(ParameterReceiver); ok {
		r.callOnParametersSet(paramReceiver, globalKey)
	}

	// Push instance onto rendering stack before calling Render
	r.renderingStack = append(r.renderingStack, instance)
	vnode := instance.Render(r)
	// Pop from rendering stack after Render completes
	r.renderingStack = r.renderingStack[:len(r.renderingStack)-1]

	return vnode
}

// cleanupUnmountedComponents removes components that are no longer in the tree
// and calls their OnUnmount lifecycle method if they implement the Unmountable interface.
func (r *RendererImpl) cleanupUnmountedComponents() {
	for key, instance := range r.instances {
		// If the component wasn't marked as active in this render, it's been unmounted
		if !r.activeKeys[key] {
			// Call OnUnmount if the component implements Unmountable
			if unmountable, ok := instance.(Unmountable); ok {
				r.callOnUnmount(unmountable, key)
			}

			// Remove from tracking maps
			delete(r.instances, key)
			delete(r.initialized, key)
		}
	}

	// Reset activeKeys for next render cycle
	r.activeKeys = make(map[string]bool)
}

// ReRender patches the DOM with minimal changes.
// This method is thread-safe and can be called from multiple goroutines.
// If multiple goroutines call this simultaneously, only one will execute at a time.
func (r *RendererImpl) ReRender() {
	// Note: RenderRoot already acquires the mutex, so this is thread-safe
	r.RenderRoot()
}

// ReRenderSlot patches only the BodyContent slot of a layout,
// preserving the layout instance and its state.
// Works by diffing the entire parent layout VDOM; only changed content is patched.
// Called when a page component (inside a layout's slot) calls StateHasChanged().
func (r *RendererImpl) ReRenderSlot(slotParent Component) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if slotParent == nil {
		return fmt.Errorf("slotParent is nil")
	}

	// 1. Get parent layout's previous VDOM from cache
	prevParentVDOM := r.instanceVDOMCache[slotParent]
	if prevParentVDOM == nil {
		// Parent not yet rendered; fall back to full re-render
		return r.reRenderFull(slotParent)
	}

	// 2. Re-render the parent layout
	// Its BodyContent field has been updated by the caller (router or child)
	// CRITICAL: Push slotParent onto rendering stack to maintain consistent key generation
	r.renderingStack = append(r.renderingStack, slotParent)
	newParentVDOM := slotParent.Render(r)
	// Pop from rendering stack after Render completes
	r.renderingStack = r.renderingStack[:len(r.renderingStack)-1]

	if newParentVDOM == nil {
		return fmt.Errorf("slotParent.Render() returned nil")
	}

	// 3. Diff the entire parent layout's VDOM and patch
	// The layout's template includes the slot content, so changes are captured
	// vdom.Patch handles the diffing and patching automatically
	vdom.Patch(r.mountID, prevParentVDOM, newParentVDOM)

	// 4. Cache the new parent VDOM for next diff
	r.instanceVDOMCache[slotParent] = newParentVDOM

	return nil
}

// reRenderFull is a helper to do a complete re-render when needed
func (r *RendererImpl) reRenderFull(component Component) error {
	newVDOM := component.Render(r)
	if newVDOM == nil {
		return fmt.Errorf("component.Render() returned nil")
	}

	vdom.RenderToSelector(r.mountID, newVDOM)
	r.instanceVDOMCache[component] = newVDOM
	r.prevVDOM = newVDOM
	return nil
}

// Navigate implements the Navigator interface.
// It delegates to the NavigationManager (router) to perform client-side navigation.
// Returns an error if no router is configured.
func (r *RendererImpl) Navigate(path string) error {
	if r.navManager == nil {
		return fmt.Errorf("no router configured for navigation")
	}
	return r.navManager.Navigate(path)
}
