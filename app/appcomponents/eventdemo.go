//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"fmt"

	"github.com/vcrobe/nojs/console"
	"github.com/vcrobe/nojs/events"
	"github.com/vcrobe/nojs/runtime"
)

// EventDemo demonstrates all supported event types with compile-time type safety.
type EventDemo struct {
	runtime.ComponentBase

	// State for tracking event interactions
	ClickCount      int
	InputValue      string
	ChangeValue     string
	LastKey         string
	LastKeyCode     string
	MousePosition   string
	IsFocused       bool
	FormSubmitted   bool
	SubmissionCount int
	KeyPressCount   int
	KeyModifiers    string
	MouseButtonInfo string
	LastEventType   string
}

// === No-Argument Events (@onclick) ===

func (c *EventDemo) HandleClick() {
	c.ClickCount++
	c.LastEventType = "Button Click"
	console.Log("Button clicked! Count:", c.ClickCount)
	c.StateHasChanged()
}

func (c *EventDemo) ResetAll() {
	c.ClickCount = 0
	c.InputValue = ""
	c.ChangeValue = ""
	c.LastKey = ""
	c.LastKeyCode = ""
	c.MousePosition = ""
	c.IsFocused = false
	c.FormSubmitted = false
	c.SubmissionCount = 0
	c.KeyPressCount = 0
	c.KeyModifiers = ""
	c.MouseButtonInfo = ""
	c.LastEventType = ""
	console.Log("All state reset!")
	c.StateHasChanged()
}

// === Change Events (@oninput, @onchange) ===

func (c *EventDemo) HandleInput(e events.ChangeEventArgs) {
	c.InputValue = e.Value
	c.LastEventType = "Input"
	console.Log("Input changed:", e.Value)
	c.StateHasChanged()
}

func (c *EventDemo) HandleChange(e events.ChangeEventArgs) {
	c.ChangeValue = e.Value
	c.LastEventType = "Change"
	console.Log("Select changed:", e.Value)
	c.StateHasChanged()
}

func (c *EventDemo) HandleTextAreaInput(e events.ChangeEventArgs) {
	c.InputValue = e.Value
	c.LastEventType = "TextArea Input"
	console.Log("TextArea input:", e.Value)
	c.StateHasChanged()
}

// === Keyboard Events (@onkeydown, @onkeyup, @onkeypress) ===

func (c *EventDemo) HandleKeyDown(e events.KeyboardEventArgs) {
	c.LastKey = e.Key
	c.LastKeyCode = e.Code
	c.LastEventType = "KeyDown"

	// Build modifier string
	mods := ""
	if e.CtrlKey {
		mods += "Ctrl+"
	}
	if e.AltKey {
		mods += "Alt+"
	}
	if e.ShiftKey {
		mods += "Shift+"
	}
	if e.MetaKey {
		mods += "Meta+"
	}
	c.KeyModifiers = mods

	console.Log(fmt.Sprintf("KeyDown: %s%s (code: %s)", mods, e.Key, e.Code))
	c.StateHasChanged()
}

func (c *EventDemo) HandleKeyUp(e events.KeyboardEventArgs) {
	c.LastEventType = "KeyUp"
	console.Log("KeyUp:", e.Key)
	c.StateHasChanged()
}

func (c *EventDemo) HandleKeyPress(e events.KeyboardEventArgs) {
	c.KeyPressCount++
	c.LastEventType = "KeyPress"
	console.Log("KeyPress:", e.Key, "Total:", c.KeyPressCount)
	c.StateHasChanged()
}

// === Focus Events (@onfocus, @onblur) ===

func (c *EventDemo) HandleFocus(e events.FocusEventArgs) {
	c.IsFocused = true
	c.LastEventType = "Focus"
	console.Log("Input focused")
	c.StateHasChanged()
}

func (c *EventDemo) HandleBlur(e events.FocusEventArgs) {
	c.IsFocused = false
	c.LastEventType = "Blur"
	console.Log("Input blurred")
	c.StateHasChanged()
}

// === Form Events (@onsubmit) ===

func (c *EventDemo) HandleSubmit(e events.FormEventArgs) {
	c.FormSubmitted = true
	c.SubmissionCount++
	c.LastEventType = "Form Submit"
	console.Log("Form submitted! Count:", c.SubmissionCount)
	c.StateHasChanged()
}

// === Mouse Events (@onmousedown, @onmouseup, @onmousemove) ===

func (c *EventDemo) HandleMouseDown(e events.MouseEventArgs) {
	buttonName := "Unknown"
	switch e.Button {
	case 0:
		buttonName = "Left"
	case 1:
		buttonName = "Middle"
	case 2:
		buttonName = "Right"
	}

	c.MouseButtonInfo = fmt.Sprintf("%s button at (%d, %d)", buttonName, e.ClientX, e.ClientY)
	c.LastEventType = "MouseDown"
	console.Log("MouseDown:", c.MouseButtonInfo)
	c.StateHasChanged()
}

func (c *EventDemo) HandleMouseUp(e events.MouseEventArgs) {
	c.LastEventType = "MouseUp"
	console.Log("MouseUp at:", e.ClientX, e.ClientY)
	c.StateHasChanged()
}

func (c *EventDemo) HandleMouseMove(e events.MouseEventArgs) {
	c.MousePosition = fmt.Sprintf("X: %d, Y: %d", e.ClientX, e.ClientY)
	c.LastEventType = "MouseMove"
	// Don't log every mousemove - too verbose
	c.StateHasChanged()
}
