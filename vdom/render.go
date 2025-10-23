//go:build js || wasm
// +build js wasm

package vdom

import (
	"syscall/js"

	"github.com/vcrobe/nojs/console"
)

func Clear(selector string) {
	if selector == "" {
		return
	}

	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}

	mount := doc.Call("querySelector", selector)
	if !mount.Truthy() {
		console.Error("Mount element not found for selector:", selector)
		return
	}

	// Set innerHTML to an empty string to clear all children.
	mount.Set("innerHTML", "")
}

// RenderToSelector mounts the VNode under the first element matching the CSS selector.
func RenderToSelector(selector string, n *VNode) {
	if n == nil || selector == "" {
		return
	}

	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}

	mount := doc.Call("querySelector", selector)

	if !mount.Truthy() {
		console.Error("Mount element not found for selector:", selector)
		return
	}

	RenderTo(mount, n)
}

// RenderTo appends the rendered node to a specific mount element.
func RenderTo(mount js.Value, n *VNode) {
	if n == nil {
		return
	}

	el := createElement(n)

	if el.Truthy() {
		mount.Call("appendChild", el)
	}
}

// setAttributeValue sets an attribute on an element, handling boolean attributes and event handlers correctly.
func setAttributeValue(el js.Value, key string, value any) {
	// Handle boolean attributes
	if boolVal, ok := value.(bool); ok {
		if boolVal {
			// For boolean attributes, set them without a value (or with empty string)
			el.Call("setAttribute", key, "")
		}
		// If false, don't set the attribute at all
		return
	}

	// Handle event handlers (functions that accept js.Value)
	if _, ok := value.(func(js.Value)); ok {
		// Event handlers should be attached via addEventListener, not setAttribute
		// Skip them here - they'll be handled separately
		return
	}

	// For all other types, convert to string and set normally
	el.Call("setAttribute", key, value)
}

// attachEventListeners processes attributes and attaches event listeners for event handlers.
// Event attributes start with "on" (e.g., onClick, onInput, onMousedown).
func attachEventListeners(el js.Value, attributes map[string]any) {
	if attributes == nil {
		return
	}

	for key, value := range attributes {
		// Check if this is an event handler (starts with "on")
		if len(key) > 2 && key[0] == 'o' && key[1] == 'n' {
			if handler, ok := value.(func(js.Value)); ok {
				// Convert "onClick" -> "click", "onInput" -> "input", etc.
				// Lowercase the first character after "on" if it's uppercase
				eventName := key[2:]
				if eventName[0] >= 'A' && eventName[0] <= 'Z' {
					eventName = string(eventName[0]+('a'-'A')) + eventName[1:]
				}

				// Wrap the handler in js.FuncOf
				cb := js.FuncOf(func(this js.Value, args []js.Value) any {
					if len(args) > 0 {
						handler(args[0])
					}
					return nil
				})

				el.Call("addEventListener", eventName, cb)
				// TODO: Store cb somewhere to release later if needed
			}
		}
	}
}

func createElement(n *VNode) js.Value {
	doc := js.Global().Get("document")
	if !doc.Truthy() || n == nil {
		return js.Undefined()
	}

	switch n.Tag {
	case "p":
		el := doc.Call("createElement", "p")

		if n.Content != "" {
			el.Set("textContent", n.Content)
		}

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		// children ignored for now
		return el
	case "div":
		el := doc.Call("createElement", "div")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Content != "" {
			el.Set("textContent", n.Content)
		}

		if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		return el
	case "input":
		el := doc.Call("createElement", "input")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		// For text input, set value if provided in Content
		if n.Content != "" {
			el.Set("value", n.Content)
		}

		return el
	case "button":
		el := doc.Call("createElement", "button")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Content != "" {
			el.Set("textContent", n.Content)
		} else if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		// Attach Go OnClick handler if present (legacy support)
		if n.OnClick != nil {
			cb := js.FuncOf(func(this js.Value, args []js.Value) any {
				n.OnClick()
				return nil
			})
			el.Call("addEventListener", "click", cb)
			// Optionally store cb somewhere to release later if needed
		}

		return el

	case "h1", "h2", "h3", "h4", "h5", "h6":
		// Handle heading tags
		el := doc.Call("createElement", n.Tag)

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Content != "" {
			el.Set("textContent", n.Content)
		}

		if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		return el

	case "ul", "ol":
		// Handle list container tags
		el := doc.Call("createElement", n.Tag)

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		return el

	case "li":
		// Handle list item tags
		el := doc.Call("createElement", "li")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Content != "" {
			el.Set("textContent", n.Content)
		}

		if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		return el

	case "select":
		// Handle select dropdown element
		el := doc.Call("createElement", "select")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		return el

	case "option":
		// Handle option element
		el := doc.Call("createElement", "option")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Content != "" {
			el.Set("textContent", n.Content)
		}

		return el

	case "textarea":
		// Handle textarea element
		el := doc.Call("createElement", "textarea")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Content != "" {
			el.Set("value", n.Content)
		}

		return el

	case "form":
		// Handle form element
		el := doc.Call("createElement", "form")

		if n.Attributes != nil {
			for k, v := range n.Attributes {
				setAttributeValue(el, k, v)
			}
			attachEventListeners(el, n.Attributes)
		}

		if n.Children != nil {
			for _, child := range n.Children {
				childEl := createElement(child)
				if childEl.Truthy() {
					el.Call("appendChild", childEl)
				}
			}
		}

		return el

	default:
		console.Error("Unsupported tag: ", n.Tag)
		return js.Undefined()
	}
}

// Patch updates the DOM by comparing old and new VDOM trees and applying minimal changes.
func Patch(mountSelector string, oldVNode, newVNode *VNode) {
	if oldVNode == nil || newVNode == nil {
		return
	}

	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}

	mount := doc.Call("querySelector", mountSelector)
	if !mount.Truthy() {
		console.Error("Mount element not found for selector:", mountSelector)
		return
	}

	// Get the root DOM element (first child of mount point)
	rootElement := mount.Get("firstChild")
	if !rootElement.Truthy() {
		// No existing DOM, just render fresh
		RenderToSelector(mountSelector, newVNode)
		return
	}

	// Patch the root element
	patchElement(rootElement, oldVNode, newVNode)
}

// patchElement updates a single DOM element based on VDOM differences.
func patchElement(domElement js.Value, oldVNode, newVNode *VNode) {
	if !domElement.Truthy() || oldVNode == nil || newVNode == nil {
		return
	}

	// If tags are different, replace the entire element
	if oldVNode.Tag != newVNode.Tag {
		newElement := createElement(newVNode)
		if newElement.Truthy() {
			parent := domElement.Get("parentNode")
			if parent.Truthy() {
				parent.Call("replaceChild", newElement, domElement)
			}
		}
		return
	}

	// Same tag - update attributes
	patchAttributes(domElement, oldVNode.Attributes, newVNode.Attributes)

	// Update content for input/textarea elements
	if newVNode.Tag == "input" || newVNode.Tag == "textarea" {
		// Only update value if element is NOT currently focused
		// This preserves the user's typing experience
		isFocused := domElement.Call("matches", ":focus")
		if !isFocused.Bool() && newVNode.Content != "" {
			currentValue := domElement.Get("value").String()
			if currentValue != newVNode.Content {
				domElement.Set("value", newVNode.Content)
			}
		}
	} else if newVNode.Tag == "select" {
		// For select elements, update the selected value
		if newVNode.Content != "" {
			domElement.Set("value", newVNode.Content)
		}
	} else {
		// Update text content for other elements
		if oldVNode.Content != newVNode.Content {
			domElement.Set("textContent", newVNode.Content)
		}
	}

	// Patch children
	patchChildren(domElement, oldVNode.Children, newVNode.Children)
}

// patchAttributes updates the attributes of a DOM element.
func patchAttributes(domElement js.Value, oldAttrs, newAttrs map[string]any) {
	// Remove old attributes that are not in new attributes
	if oldAttrs != nil {
		for key := range oldAttrs {
			if _, exists := newAttrs[key]; !exists {
				// Skip event handlers (they start with "on")
				if len(key) > 2 && key[0] == 'o' && key[1] == 'n' {
					continue
				}
				domElement.Call("removeAttribute", key)
			}
		}
	}

	// Set new attributes
	if newAttrs != nil {
		for key, value := range newAttrs {
			// Skip event handlers - they're attached separately
			if len(key) > 2 && key[0] == 'o' && key[1] == 'n' {
				continue
			}

			// Check if attribute changed
			if oldAttrs == nil || oldAttrs[key] != value {
				setAttributeValue(domElement, key, value)
			}
		}
	}
}

// patchChildren updates the children of a DOM element.
func patchChildren(domElement js.Value, oldChildren, newChildren []*VNode) {
	oldLen := len(oldChildren)
	newLen := len(newChildren)
	minLen := oldLen
	if newLen < minLen {
		minLen = newLen
	}

	// Get the DOM children
	domChildren := domElement.Get("childNodes")

	// Patch existing children
	for i := 0; i < minLen; i++ {
		childElement := domChildren.Call("item", i)
		if childElement.Truthy() {
			patchElement(childElement, oldChildren[i], newChildren[i])
		}
	}

	// Add new children if newChildren is longer
	if newLen > oldLen {
		for i := oldLen; i < newLen; i++ {
			newChild := createElement(newChildren[i])
			if newChild.Truthy() {
				domElement.Call("appendChild", newChild)
			}
		}
	}

	// Remove extra children if oldChildren is longer
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			childElement := domChildren.Call("item", i)
			if childElement.Truthy() {
				domElement.Call("removeChild", childElement)
			}
		}
	}
}
