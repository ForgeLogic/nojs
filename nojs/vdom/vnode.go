package vdom

// VNode represents a virtual DOM node.
// This core file has NO build tags, making it available to both WASM and native test builds.
type VNode struct {
	Tag            string         // The HTML tag name
	Attributes     map[string]any // The attributes of the node
	Children       []*VNode       // The child nodes
	Content        string         // The content of the node
	OnClick        func()         // Optional click event handler
	Key            any            // Optional key for list reconciliation (used in {@for} loops)
	ComponentKey   string         // Key for component-level reconciliation (used in router navigation)
	eventCallbacks []any          // Stores js.Func objects for cleanup (interface{} to avoid build tag issues)
}

// NewVNode creates a new VNode.
func NewVNode(tag string, attributes map[string]any, children []*VNode, content string) *VNode {
	var onClick func()
	if attributes != nil {
		if v, ok := attributes["onClick"]; ok {
			if f, ok := v.(func()); ok {
				onClick = f
				// Remove from attributes so it doesn't get rendered as an HTML attribute
				delete(attributes, "onClick")
			}
		}
	}
	return &VNode{
		Tag:        tag,
		Attributes: attributes,
		Children:   children,
		Content:    content,
		OnClick:    onClick,
	}
}

// SetContent updates the Content field of the VNode.
func (v *VNode) SetContent(content string) {
	v.Content = content
}

// Text creates a pure text node VNode (no HTML element wrapper).
// This uses the special "#text" tag which renders as document.createTextNode() in the browser.
// Use this for creating text content without any surrounding HTML element.
func Text(content string) *VNode {
	return &VNode{
		Tag:     "#text",
		Content: content,
	}
}

// Paragraph creates a <p> VNode with the given text as its child and allows passing attributes.
func Paragraph(text string, attrs map[string]any) *VNode {
	return NewVNode("p", attrs, nil, text)
}

// InputText returns a VNode representing an <input type="text"> element.
// Optionally accepts a map of attributes (e.g., {"placeholder": "Type here"}).
func InputText(attrs map[string]any) *VNode {
	if attrs == nil {
		attrs = make(map[string]any)
	}
	attrs["type"] = "text"
	return NewVNode("input", attrs, nil, "")
}

// Div creates a <div> VNode with the given children and allows passing attributes.
func Div(attrs map[string]any, children ...*VNode) *VNode {
	return NewVNode("div", attrs, children, "")
}

// Button creates a <button> VNode with the given children and allows passing attributes.
func Button(content string, attrs map[string]any, children ...*VNode) *VNode {
	return NewVNode("button", attrs, children, content)
}

// AddEventCallback stores a js.Func (as interface{}) for later cleanup.
// This is called from WASM-only code in render.go.
func (v *VNode) AddEventCallback(cb any) {
	v.eventCallbacks = append(v.eventCallbacks, cb)
}

// GetEventCallbacks returns all stored event callbacks.
// This is used by WASM-only code to access the callbacks for cleanup.
func (v *VNode) GetEventCallbacks() []any {
	return v.eventCallbacks
}

// ClearEventCallbacks clears the event callbacks slice without releasing them.
// The actual release is done in WASM-only code in render.go.
func (v *VNode) ClearEventCallbacks() {
	v.eventCallbacks = nil
}
