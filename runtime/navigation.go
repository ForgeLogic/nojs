//go:build js && wasm

package runtime

// NavigationManager defines the contract for a client-side router
// that integrates with the No-JS framework's rendering lifecycle.
//
// The framework's core is router-agnostic. Any router implementation
// that fulfills this interface can be plugged into the Renderer.
type NavigationManager interface {
	// Start initializes the router. It must read the initial browser URL,
	// determine the initial component, call the onChange callback with it,
	// and begin listening for browser history events (like popstate).
	//
	// The onChange callback is provided by the application (usually in main.go)
	// and is responsible for updating the renderer's current component and
	// triggering a re-render.
	Start(onChange func(newComponent Component)) error

	// Navigate programmatically changes the browser URL using the History API
	// (or hash navigation) and triggers the onChange callback with the component
	// for the new path.
	//
	// This is called when components request navigation (e.g., via a Link component).
	Navigate(path string) error

	// GetComponentForPath resolves a URL path to its corresponding
	// component instance based on the router's configuration.
	//
	// Returns (component, true) if a route matches, or (nil, false) if not found.
	GetComponentForPath(path string) (Component, bool)
}

// Navigator is an interface that components can use to request navigation.
// This is provided to components via the runtime.ComponentBase.
//
// Components should call Navigate(path) to trigger client-side routing,
// which will update the browser URL and render the new component without
// a full page reload.
type Navigator interface {
	// Navigate requests navigation to a new path.
	// Returns an error if navigation fails.
	Navigate(path string) error
}
