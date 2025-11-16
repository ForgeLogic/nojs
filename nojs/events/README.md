# Events Package

The `events` package provides type-safe event handling for the nojs framework. It enables compile-time validation of event handler signatures and bridges Go methods to browser events.

## Overview

This package contains three main components:

1. **Event Argument Types** (`events.go`) - Strongly-typed structs for event data
2. **Event Registry** (`registry.go`) - Maps events to their expected signatures for validation
3. **Event Adapters** (`adapters.go`) - Bridges Go handlers to JavaScript events

## Event Types

### ChangeEventArgs
Used for: `@oninput`, `@onchange`  
Supported elements: `<input>`, `<textarea>`, `<select>`

```go
func (c *MyComponent) HandleInput(e events.ChangeEventArgs) {
    c.Value = e.Value
    c.StateHasChanged()
}
```

### KeyboardEventArgs
Used for: `@onkeydown`, `@onkeyup`, `@onkeypress`  
Supported elements: `<input>`, `<textarea>`, `<div>`

```go
func (c *MyComponent) HandleKey(e events.KeyboardEventArgs) {
    if e.Key == "Enter" && e.CtrlKey {
        c.Submit()
    }
    c.StateHasChanged()
}
```

### MouseEventArgs
Used for: `@onmousedown`, `@onmouseup`, `@onmousemove`  
Supported elements: `<button>`, `<div>`, `<span>`, `<img>`, `<a>`, `<canvas>`

```go
func (c *MyComponent) HandleMouseMove(e events.MouseEventArgs) {
    c.CursorX = e.ClientX
    c.CursorY = e.ClientY
    c.StateHasChanged()
}
```

### FocusEventArgs
Used for: `@onfocus`, `@onblur`  
Supported elements: `<input>`, `<textarea>`, `<select>`, `<button>`

```go
func (c *MyComponent) HandleFocus(e events.FocusEventArgs) {
    c.IsFocused = true
    c.StateHasChanged()
}
```

### FormEventArgs
Used for: `@onsubmit`  
Supported elements: `<form>`

```go
func (c *MyComponent) HandleSubmit(e events.FormEventArgs) {
    // Form submission is automatically prevented
    c.ProcessForm()
    c.StateHasChanged()
}
```

## No-Argument Events

Some events don't require arguments:

### onclick
Used for: `@onclick`  
Supported elements: `<button>`, `<a>`, `<div>`, `<span>`, `<p>`, `<img>`

```go
func (c *MyComponent) HandleClick() {
    c.Counter++
    c.StateHasChanged()
}
```

## Compile-Time Validation

The compiler validates event handlers at build time:

```html
<!-- ✅ Valid: Handler signature matches event -->
<button @onclick="HandleClick">Click</button>
<input @oninput="HandleInput" />

<!-- ❌ Invalid: Signature mismatch -->
<button @onclick="HandleInput">Click</button>
<!-- Error: Handler 'HandleInput' has signature 'func(events.ChangeEventArgs)' 
     but @onclick on <button> requires 'func()' -->
```

## Supported Events by Phase

### Phase 1 (MVP)
- ✅ `@onclick` (no args)
- ✅ `@oninput` (ChangeEventArgs)
- ✅ `@onchange` (ChangeEventArgs)

### Phase 2
- ✅ `@onkeydown`, `@onkeyup`, `@onkeypress` (KeyboardEventArgs)
- ✅ `@onfocus`, `@onblur` (FocusEventArgs)
- ✅ `@onsubmit` (FormEventArgs)

### Phase 3
- ✅ `@onmousedown`, `@onmouseup`, `@onmousemove` (MouseEventArgs)

## Implementation Notes

- All event files use `//go:build js && wasm` build tag
- Adapters automatically extract event properties from `syscall/js.Value`
- Form submissions automatically call `preventDefault()`
- Event validation happens at compile time, not runtime
