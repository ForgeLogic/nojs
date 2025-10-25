//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
)

type LifecycleDemo struct {
	runtime.ComponentBase

	IsLoading bool
	UserId    int  `nojs:"state"` // Internal state - won't be copied by ApplyProps
	ShowChild bool `nojs:"state"` // Control whether child is mounted
}

var _ runtime.Initializer = (*LifecycleDemo)(nil)

func (c *LifecycleDemo) OnInit() {
	c.IsLoading = true
	c.ShowChild = true // Start with child visible
}

func (c *LifecycleDemo) User1Click() {
	c.UserId = 1
	c.StateHasChanged()
}

func (c *LifecycleDemo) User2Click() {
	c.UserId = 2
	c.StateHasChanged()
}

func (c *LifecycleDemo) ToggleChild() {
	c.ShowChild = !c.ShowChild
	c.StateHasChanged()
}
