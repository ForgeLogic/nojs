//go:build js && wasm

package events

import "syscall/js"

// AdaptClickEvent creates a JavaScript-compatible event handler from a Go handler
// that expects ClickEventArgs. This is used for @onclick events with event arguments.
func AdaptClickEvent(handler func(ClickEventArgs)) func(js.Value) {
	return func(e js.Value) {
		args := ClickEventArgs{
			EventBase: NewEventBase(e),
			ClientX:   e.Get("clientX").Int(),
			ClientY:   e.Get("clientY").Int(),
			Button:    e.Get("button").Int(),
			AltKey:    e.Get("altKey").Bool(),
			CtrlKey:   e.Get("ctrlKey").Bool(),
			ShiftKey:  e.Get("shiftKey").Bool(),
			MetaKey:   e.Get("metaKey").Bool(),
		}
		handler(args)
	}
}

// AdaptChangeEvent creates a JavaScript-compatible event handler from a Go handler
// that expects ChangeEventArgs. This is used for @oninput and @onchange events.
func AdaptChangeEvent(handler func(ChangeEventArgs)) func(js.Value) {
	return func(e js.Value) {
		args := ChangeEventArgs{
			EventBase: NewEventBase(e),
			Value:     e.Get("target").Get("value").String(),
		}
		handler(args)
	}
}

// AdaptKeyboardEvent creates a JavaScript-compatible event handler from a Go handler
// that expects KeyboardEventArgs. This is used for @onkeydown, @onkeyup, @onkeypress events.
func AdaptKeyboardEvent(handler func(KeyboardEventArgs)) func(js.Value) {
	return func(e js.Value) {
		args := KeyboardEventArgs{
			EventBase: NewEventBase(e),
			Key:       e.Get("key").String(),
			Code:      e.Get("code").String(),
			AltKey:    e.Get("altKey").Bool(),
			CtrlKey:   e.Get("ctrlKey").Bool(),
			ShiftKey:  e.Get("shiftKey").Bool(),
			MetaKey:   e.Get("metaKey").Bool(),
		}
		handler(args)
	}
}

// AdaptMouseEvent creates a JavaScript-compatible event handler from a Go handler
// that expects MouseEventArgs. This is used for @onmousedown, @onmouseup, @onmousemove events.
func AdaptMouseEvent(handler func(MouseEventArgs)) func(js.Value) {
	return func(e js.Value) {
		args := MouseEventArgs{
			EventBase: NewEventBase(e),
			ClientX:   e.Get("clientX").Int(),
			ClientY:   e.Get("clientY").Int(),
			Button:    e.Get("button").Int(),
			AltKey:    e.Get("altKey").Bool(),
			CtrlKey:   e.Get("ctrlKey").Bool(),
			ShiftKey:  e.Get("shiftKey").Bool(),
			MetaKey:   e.Get("metaKey").Bool(),
		}
		handler(args)
	}
}

// AdaptFocusEvent creates a JavaScript-compatible event handler from a Go handler
// that expects FocusEventArgs. This is used for @onfocus and @onblur events.
func AdaptFocusEvent(handler func(FocusEventArgs)) func(js.Value) {
	return func(e js.Value) {
		args := FocusEventArgs{
			EventBase: NewEventBase(e),
		}
		handler(args)
	}
}

// AdaptFormEvent creates a JavaScript-compatible event handler from a Go handler
// that expects FormEventArgs. This is used for @onsubmit events.
func AdaptFormEvent(handler func(FormEventArgs)) func(js.Value) {
	return func(e js.Value) {
		args := FormEventArgs{
			EventBase: NewEventBase(e),
		}
		handler(args)
	}
}

// AdaptNoArgEvent creates a JavaScript-compatible event handler from a Go handler
// that expects no arguments. This is used for @onclick events.
func AdaptNoArgEvent(handler func()) func(js.Value) {
	return func(e js.Value) {
		handler()
	}
}
