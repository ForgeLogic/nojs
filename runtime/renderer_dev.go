//go:build (js || wasm) && dev
// +build js wasm
// +build dev

package runtime

// callOnMount invokes the OnMount lifecycle method in development mode.
// In dev mode, panics propagate to aid debugging and fast failure.
func (r *RendererImpl) callOnMount(mountable Mountable, key string) {
	mountable.OnMount()
}

// callOnParametersSet invokes the OnParametersSet lifecycle method in development mode.
// In dev mode, panics propagate to aid debugging and fast failure.
func (r *RendererImpl) callOnParametersSet(receiver ParameterReceiver, key string) {
	receiver.OnParametersSet()
}

// callOnUnmount invokes the OnUnmount lifecycle method in development mode.
// In dev mode, panics propagate to aid debugging and fast failure.
func (r *RendererImpl) callOnUnmount(unmountable Unmountable, key string) {
	unmountable.OnUnmount()
}
