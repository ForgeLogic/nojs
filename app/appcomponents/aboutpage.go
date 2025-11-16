//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
)

// AboutPage is the component rendered for the "/about" route.
type AboutPage struct {
	runtime.ComponentBase
}
