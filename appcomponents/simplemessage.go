//go:build js || wasm
// +build js wasm

package appcomponents

type SimpleMessage struct {
	FirstProp string
}

func (qwe *SimpleMessage) Increment() {
	println("Clicked")
}
