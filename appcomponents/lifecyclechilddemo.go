//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"context"
	"time"

	"github.com/vcrobe/nojs/runtime"
)

type LifecycleChildDemo struct {
	runtime.ComponentBase

	UserId int

	// Internal state for timer demo
	ctx       context.Context    `nojs:"state"`
	cancel    context.CancelFunc `nojs:"state"`
	TickCount int                `nojs:"state"`
}

var _ runtime.Initializer = (*LifecycleChildDemo)(nil)
var _ runtime.ParameterReceiver = (*LifecycleChildDemo)(nil)
var _ runtime.Cleaner = (*LifecycleChildDemo)(nil)

func (c *LifecycleChildDemo) OnInit() {
	println("calling LcChildDemo.OnInit - starting timer for UserId:", c.UserId)

	// Start a timer to demonstrate cleanup
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.runTimer()
}

func (c *LifecycleChildDemo) OnParametersSet() {
}

func (c *LifecycleChildDemo) runTimer() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			println("Timer stopped for UserId:", c.UserId, "- cleanup complete")
			return
		case <-ticker.C:
			c.TickCount++
			c.StateHasChanged()
		}
	}
}

func (c *LifecycleChildDemo) OnDestroy() {
	println("calling LcChildDemo.OnDestroy for UserId:", c.UserId, "- stopping timer")
	if c.cancel != nil {
		c.cancel()
	}
}
