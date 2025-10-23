//go:build js && wasm

package events

// ChangeEventArgs represents the data passed from input/select/textarea change events.
// This struct provides type-safe access to the current value of form elements.
type ChangeEventArgs struct {
	// Value is the current value of the input element.
	// For text inputs, this is the text content.
	// For select elements, this is the selected option's value.
	// For checkboxes, this will be "true" or "false".
	Value string
}

// KeyboardEventArgs represents the data passed from keyboard events.
// Used for @onkeydown, @onkeyup, @onkeypress handlers.
type KeyboardEventArgs struct {
	Key      string // The key value of the key pressed (e.g., "a", "Enter", "Escape")
	Code     string // The physical key code (e.g., "KeyA", "Enter")
	AltKey   bool   // Whether the Alt key was pressed
	CtrlKey  bool   // Whether the Ctrl key was pressed
	ShiftKey bool   // Whether the Shift key was pressed
	MetaKey  bool   // Whether the Meta (Command/Windows) key was pressed
}

// MouseEventArgs represents the data passed from mouse events.
// Used for @onmousedown, @onmouseup, @onmousemove handlers.
type MouseEventArgs struct {
	ClientX  int  // X coordinate relative to the viewport
	ClientY  int  // Y coordinate relative to the viewport
	Button   int  // Which mouse button was pressed (0=left, 1=middle, 2=right)
	AltKey   bool // Whether the Alt key was pressed
	CtrlKey  bool // Whether the Ctrl key was pressed
	ShiftKey bool // Whether the Shift key was pressed
	MetaKey  bool // Whether the Meta key was pressed
}

// FocusEventArgs represents the data passed from focus/blur events.
// Currently minimal, can be extended if needed.
type FocusEventArgs struct {
	// Reserved for future use
}

// FormEventArgs represents the data passed from form submission events.
// Used for @onsubmit handlers.
type FormEventArgs struct {
	// Reserved for future use
	// May include form data extraction in the future
}
