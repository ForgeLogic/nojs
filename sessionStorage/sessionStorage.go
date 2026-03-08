//go:build js || wasm

package sessionStorage

import (
	"syscall/js"
)

func Clear() {
	storage := js.Global().Get("sessionStorage")
	storage.Call("clear")
}

func GetItem(keyName string) string {
	storage := js.Global().Get("sessionStorage")
	keyValue := storage.Call("getItem", keyName)
	return keyValue.String()
}

func Length() int {
	storage := js.Global().Get("sessionStorage")
	len := storage.Get("length")
	return len.Int()
}

func SetItem(keyName, keyValue string) {
	storage := js.Global().Get("sessionStorage")
	storage.Call("setItem", keyName, keyValue)
}

func RemoveItem(keyName string) {
	storage := js.Global().Get("sessionStorage")
	storage.Call("removeItem", keyName)
}
