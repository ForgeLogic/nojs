//go:build js || wasm
// +build js wasm

package main

import (
	"syscall/js"

	"github.com/vcrobe/nojs/dialogs"
	"github.com/vcrobe/nojs/vdom"
)

func add(this js.Value, args []js.Value) interface{} {
	a := args[0].Int()
	b := args[1].Int()
	return js.ValueOf(a + b)
}

func showPrompt() {
	name := dialogs.Prompt("write your name")

	if name == "<null>" {
		println("you pressed the cancel button")
	} else if name == "" {
		println("the string is empty")
	} else {
		println("your name is", name)
	}
}

func callJsFunction() {
	js.Global().Call("calledFromGoWasm", "Hello from Go!")
}

func testButtonClick() {
	println("Button was clicked!")
}

func main() {
	// Export the `add` function to JavaScript
	js.Global().Set("add", js.FuncOf(add))

	// Call the JavaScript function
	callJsFunction()

	div := vdom.Div(map[string]any{"id": "test-div", "data-attr": -1.3},
		vdom.Paragraph("Paragraph with attributes", map[string]any{"id": "test-paragraph"}),
		vdom.Div(nil,
			vdom.Paragraph("Simple paragraph tag", nil),
		),
		vdom.Button("Click me", map[string]any{
			"onClick": func() { testButtonClick() },
		}, nil),
	)
	vdom.RenderToSelector("#app", div)

	// Keep the Go program running
	select {}
}
