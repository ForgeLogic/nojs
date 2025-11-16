//go:build js && wasm

package events

import "syscall/js"

// EventBase provides common functionality for all DOM events.
// Components can embed this to gain access to PreventDefault and StopPropagation.
type EventBase struct {
	jsEvent               js.Value
	preventDefaultCalled  bool
	stopPropagationCalled bool
}

// NewEventBase creates a new EventBase from a JavaScript event object.
// This is called by the adapter functions.
func NewEventBase(jsEvent js.Value) EventBase {
	return EventBase{
		jsEvent: jsEvent,
	}
}

// PreventDefault prevents the browser's default action for this event.
// For example, prevents form submission or link navigation.
func (e *EventBase) PreventDefault() {
	if !e.preventDefaultCalled {
		e.jsEvent.Call("preventDefault")
		e.preventDefaultCalled = true
	}
}

// StopPropagation stops the event from bubbling up the DOM tree.
func (e *EventBase) StopPropagation() {
	if !e.stopPropagationCalled {
		e.jsEvent.Call("stopPropagation")
		e.stopPropagationCalled = true
	}
}

// IsDefaultPrevented returns whether preventDefault was called.
func (e *EventBase) IsDefaultPrevented() bool {
	return e.preventDefaultCalled
}

// IsPropagationStopped returns whether stopPropagation was called.
func (e *EventBase) IsPropagationStopped() bool {
	return e.stopPropagationCalled
}

// ClickEventArgs represents the data passed from click events.
// Used for @onclick handlers that need event details.
type ClickEventArgs struct {
	EventBase
	ClientX  int  // X coordinate relative to the viewport
	ClientY  int  // Y coordinate relative to the viewport
	Button   int  // Which mouse button was pressed (0=left, 1=middle, 2=right)
	AltKey   bool // Whether the Alt key was pressed
	CtrlKey  bool // Whether the Ctrl key was pressed
	ShiftKey bool // Whether the Shift key was pressed
	MetaKey  bool // Whether the Meta key was pressed
}

// ChangeEventArgs represents the data passed from input/select/textarea change events.
// This struct provides type-safe access to the current value of form elements.
type ChangeEventArgs struct {
	EventBase
	// Value is the current value of the input element.
	// For text inputs, this is the text content.
	// For select elements, this is the selected option's value.
	// For checkboxes, this will be "true" or "false".
	Value string
}

// KeyboardEventArgs represents the data passed from keyboard events.
// Used for @onkeydown, @onkeyup, @onkeypress handlers.
type KeyboardEventArgs struct {
	EventBase
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
	EventBase
	ClientX  int  // X coordinate relative to the viewport
	ClientY  int  // Y coordinate relative to the viewport
	Button   int  // Which mouse button was pressed (0=left, 1=middle, 2=right)
	AltKey   bool // Whether the Alt key was pressed
	CtrlKey  bool // Whether the Ctrl key was pressed
	ShiftKey bool // Whether the Shift key was pressed
	MetaKey  bool // Whether the Meta key was pressed
}

// FocusEventArgs represents the data passed from focus/blur events.
type FocusEventArgs struct {
	EventBase
}

// FormEventArgs represents the data passed from form submission events.
// Used for @onsubmit handlers.
type FormEventArgs struct {
	EventBase
}
