//go:build (js || wasm) && dev
// +build js wasm
// +build dev

package runtime

// callOnInit invokes the OnInit lifecycle method in development mode.
// In dev mode, panics propagate to aid debugging and fast failure.
func (r *Renderer) callOnInit(initializer Initializer, key string) {
	initializer.OnInit()
}

// callOnParametersSet invokes the OnParametersSet lifecycle method in development mode.
// In dev mode, panics propagate to aid debugging and fast failure.
func (r *Renderer) callOnParametersSet(receiver ParameterReceiver, key string) {
	receiver.OnParametersSet()
}

// callOnDestroy invokes the OnDestroy lifecycle method in development mode.
// In dev mode, panics propagate to aid debugging and fast failure.
func (r *Renderer) callOnDestroy(cleaner Cleaner, key string) {
	cleaner.OnDestroy()
}
