//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
)

type App struct {
	runtime.ComponentBase
}

func NewApp() *App {
	return &App{}
}
