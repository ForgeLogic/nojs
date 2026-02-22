//go:build js || wasm

package pages

import (
	"github.com/ForgeLogic/nojs/runtime"
)

var routerDemoIDs = []string{"42", "go-wasm", "hello", "framework", "nojs", "2026"}

// routerParamsRenderCount persists across instance recreations so the counter
// accurately reflects total navigations to this route, not just renders of
// the current instance (which resets on every param change).
var routerParamsRenderCount int

// RouterParamsPage demonstrates URL route parameters and programmatic navigation.
type RouterParamsPage struct {
	runtime.ComponentBase

	ID          string
	RenderCount int

	nextIDIndex int
}

func (c *RouterParamsPage) OnParametersSet() {
	routerParamsRenderCount++
	c.RenderCount = routerParamsRenderCount

	println("RouterParamsPage: OnParametersSet called with ID =", c.ID)

	// Sync nextIDIndex so GoToNext() always advances from the current ID,
	// even when the instance is freshly created by the router (which resets all fields).
	for i, id := range routerDemoIDs {
		if id == c.ID {
			c.nextIDIndex = (i + 1) % len(routerDemoIDs)
			break
		}
	}
}

func (c *RouterParamsPage) GoToNext() {
	next := routerDemoIDs[c.nextIDIndex%len(routerDemoIDs)]
	c.nextIDIndex++
	c.Navigate("/demo/router/" + next)
}
