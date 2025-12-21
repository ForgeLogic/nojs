package context

// MainLayoutCtx acts as a shared state bridge.
type MainLayoutCtx struct {
	Title string
	// OnUpdate is a callback to trigger a UI refresh
	// on the component that "owns" the layout.
	OnUpdate func()
}

func (c *MainLayoutCtx) SetTitle(t string) {
	c.Title = t
	if c.OnUpdate != nil {
		c.OnUpdate()
	}
}
