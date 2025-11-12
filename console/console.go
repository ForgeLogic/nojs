//go:build js || wasm

package console

import (
	"syscall/js"
)

func Log(args ...any) {
	console := js.Global().Get("console")
	console.Call("log", args...)
}

func Warn(args ...any) {
	console := js.Global().Get("console")
	console.Call("warn", args...)
}

func Error(args ...any) {
	console := js.Global().Get("console")
	console.Call("error", args...)
}
