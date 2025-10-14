//go:build js || wasm
// +build js wasm

package appcomponents

type App struct{}

func NewApp() *App {
	return &App{}
}
