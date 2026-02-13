//go:build js || wasm
// +build js wasm

package layouts

import (
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/vdom"
)

// AdminLayout is a sublayout for admin pages.
// It wraps the admin content with a sidebar and admin-specific navigation.
// State is preserved across navigations within the /admin/* route space.
type AdminLayout struct {
	runtime.ComponentBase

	BodyContent []*vdom.VNode
	selectedNav string
}

func NewAdminLayout() runtime.Component {
	return &AdminLayout{
		selectedNav: "dashboard",
	}
}

func (a *AdminLayout) OnInit() {
	// Initialize any admin-specific state
}

func (a *AdminLayout) OnMount() {
	// Called when AdminLayout is first mounted
	println("[AdminLayout] Mounted")
}

func (a *AdminLayout) OnUnmount() {
	// Called when AdminLayout is being removed from the tree
	println("[AdminLayout] Unmounted")
}

// SetBodyContent sets the content of the BodyContent slot.
// Called by the router engine when linking layout chains.
func (a *AdminLayout) SetBodyContent(content []*vdom.VNode) {
	a.BodyContent = content
}
