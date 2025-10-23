package events

// EventSignature defines the expected handler signature for a specific event.
// The compiler uses this to validate that component methods match the required signature.
type EventSignature struct {
	// EventName is the event attribute name without @ prefix (e.g., "onclick", "oninput")
	EventName string

	// SupportedTags lists HTML elements that support this event
	SupportedTags []string

	// ExpectedSig is a human-readable signature string for error messages
	ExpectedSig string

	// RequiresArgs indicates if the handler expects event arguments
	RequiresArgs bool

	// ArgsType is the fully-qualified Go type name for event arguments
	// (e.g., "events.ChangeEventArgs")
	ArgsType string
}

// EventRegistry maps event names to their expected signatures.
// This is used by the compiler for compile-time validation.
var EventRegistry = map[string]EventSignature{
	// Phase 1: Core events (MVP)
	"onclick": {
		EventName:     "onclick",
		SupportedTags: []string{"button", "a", "div", "span", "p", "img"},
		ExpectedSig:   "func()",
		RequiresArgs:  false,
	},
	"oninput": {
		EventName:     "oninput",
		SupportedTags: []string{"input", "textarea"},
		ExpectedSig:   "func(events.ChangeEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.ChangeEventArgs",
	},
	"onchange": {
		EventName:     "onchange",
		SupportedTags: []string{"input", "select", "textarea"},
		ExpectedSig:   "func(events.ChangeEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.ChangeEventArgs",
	},

	// Phase 2: Keyboard events
	"onkeydown": {
		EventName:     "onkeydown",
		SupportedTags: []string{"input", "textarea", "div"},
		ExpectedSig:   "func(events.KeyboardEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.KeyboardEventArgs",
	},
	"onkeyup": {
		EventName:     "onkeyup",
		SupportedTags: []string{"input", "textarea", "div"},
		ExpectedSig:   "func(events.KeyboardEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.KeyboardEventArgs",
	},
	"onkeypress": {
		EventName:     "onkeypress",
		SupportedTags: []string{"input", "textarea", "div"},
		ExpectedSig:   "func(events.KeyboardEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.KeyboardEventArgs",
	},

	// Phase 2: Focus events
	"onfocus": {
		EventName:     "onfocus",
		SupportedTags: []string{"input", "textarea", "select", "button"},
		ExpectedSig:   "func(events.FocusEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.FocusEventArgs",
	},
	"onblur": {
		EventName:     "onblur",
		SupportedTags: []string{"input", "textarea", "select", "button"},
		ExpectedSig:   "func(events.FocusEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.FocusEventArgs",
	},

	// Phase 2: Form events
	"onsubmit": {
		EventName:     "onsubmit",
		SupportedTags: []string{"form"},
		ExpectedSig:   "func(events.FormEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.FormEventArgs",
	},

	// Phase 3: Mouse events
	"onmousedown": {
		EventName:     "onmousedown",
		SupportedTags: []string{"button", "div", "span", "img", "a"},
		ExpectedSig:   "func(events.MouseEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.MouseEventArgs",
	},
	"onmouseup": {
		EventName:     "onmouseup",
		SupportedTags: []string{"button", "div", "span", "img", "a"},
		ExpectedSig:   "func(events.MouseEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.MouseEventArgs",
	},
	"onmousemove": {
		EventName:     "onmousemove",
		SupportedTags: []string{"div", "span", "canvas"},
		ExpectedSig:   "func(events.MouseEventArgs)",
		RequiresArgs:  true,
		ArgsType:      "events.MouseEventArgs",
	},
}

// GetEventSignature returns the signature for an event name.
// Returns nil if the event is not registered.
func GetEventSignature(eventName string) *EventSignature {
	if sig, ok := EventRegistry[eventName]; ok {
		return &sig
	}
	return nil
}

// IsEventSupported checks if an event is supported on a specific HTML tag.
func IsEventSupported(eventName, tagName string) bool {
	sig := GetEventSignature(eventName)
	if sig == nil {
		return false
	}

	for _, supportedTag := range sig.SupportedTags {
		if supportedTag == tagName {
			return true
		}
	}
	return false
}
