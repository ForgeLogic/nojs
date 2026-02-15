package layouts

import (
	"github.com/vcrobe/app/internal/app/context"
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/vdom"
)

// RootLayout is the root layout component for the application.
type MainLayout struct {
	runtime.ComponentBase

	MainLayoutCtx *context.MainLayoutCtx
	BodyContent   []*vdom.VNode
}

func (c *MainLayout) OnMount() {
	// We ensure the callback points to this component's refresh logic
	if c.MainLayoutCtx != nil {
		c.MainLayoutCtx.OnUpdate = c.StateHasChanged
	}
}

// SetBodyContent sets the content of the BodyContent slot.
// Called by the router engine when linking layout chains.
func (c *MainLayout) SetBodyContent(content []*vdom.VNode) {
	c.BodyContent = content
}
