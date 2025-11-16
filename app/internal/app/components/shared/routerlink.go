//go:build js && wasm

package shared

import (
	"github.com/vcrobe/nojs/events"
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/vdom"
)

// RouterLink is a component for client-side navigation.
// It renders an <a> tag that uses the router for navigation without page reloads.
//
// Props:
//   - Href: The path to navigate to (e.g., "/about", "/users/123")
//   - Children: The content to display inside the link (text, other components, etc.)
//
// Example usage in a template:
//
//	<RouterLink Href="/about">
//	    <span>Go to About Page</span>
//	</RouterLink>
type RouterLink struct {
	runtime.ComponentBase

	// Href is the destination path for navigation
	Href string

	// Children contains the content projected into the link
	Children []*vdom.VNode
}

// HandleClick is called when the link is clicked.
// It prevents the default browser navigation and uses the router instead.
func (c *RouterLink) HandleClick(e events.ClickEventArgs) {
	// Prevent the browser from navigating (which would reload the page)
	e.PreventDefault()

	// Use the framework's client-side router to navigate
	if err := c.Navigate(c.Href); err != nil {
		println("[RouterLink] Navigation error:", err.Error())
	}
}
