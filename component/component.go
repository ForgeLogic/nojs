package component

import "github.com/vcrobe/nojs/vdom"

// Component interface defines the structure for all components in the framework.
type Component interface {
	Render() *vdom.VNode // Renders the component's virtual DOM
}
