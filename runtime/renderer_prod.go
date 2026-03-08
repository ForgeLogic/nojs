//go:build (js || wasm) && !dev
// +build js wasm
// +build !dev

package runtime

import "fmt"

// callOnMount invokes the OnMount lifecycle method in production mode.
// In production mode, panics are recovered and logged to prevent application crashes.
func (r *RendererImpl) callOnMount(mountable Mountable, key string) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("ERROR: OnMount panic in component %s: %v\n", key, rec)
			// In a real production environment, this could be sent to an error tracking service
		}
	}()
	mountable.OnMount()
}

// callOnParametersSet invokes the OnParametersSet lifecycle method in production mode.
// In production mode, panics are recovered and logged to prevent application crashes.
func (r *RendererImpl) callOnParametersSet(receiver ParameterReceiver, key string) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("ERROR: OnParametersSet panic in component %s: %v\n", key, rec)
			// In a real production environment, this could be sent to an error tracking service
		}
	}()
	receiver.OnParametersSet()
}

// callOnUnmount invokes the OnUnmount lifecycle method in production mode.
// In production mode, panics are recovered and logged to prevent application crashes.
func (r *RendererImpl) callOnUnmount(unmountable Unmountable, key string) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("ERROR: OnUnmount panic in component %s: %v\n", key, rec)
			// In a real production environment, this could be sent to an error tracking service
		}
	}()
	unmountable.OnUnmount()
}
