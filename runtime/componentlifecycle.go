//go:build js || wasm
// +build js wasm

package runtime

// Initializer is implemented by components that need one-time initialization.
// OnInit is called once after the component instance is created, before the first render.
//
// Example:
//
//	type UserProfile struct {
//	    runtime.ComponentBase
//	    UserID    int
//	    User      *UserData
//	    IsLoading bool
//	}
//
//	func (c *UserProfile) OnInit() {
//	    c.IsLoading = true
//	    go c.fetchUserData() // Launch async task
//	}
//
//	func (c *UserProfile) fetchUserData() {
//	    // Simulate API call...
//	    fetchedUser := &UserData{Name: "Fetched Name"}
//	    c.User = fetchedUser
//	    c.IsLoading = false
//	    c.StateHasChanged() // Trigger re-render
//	}
type Initializer interface {
	OnInit()
}

// ParameterReceiver is implemented by components that need to react to parameter changes.
// OnParametersSet is called every time the component receives (potentially new) parameters
// from its parent, just before Render is called. This includes the initial render.
//
// Example:
//
//	type DataDisplay struct {
//	    runtime.ComponentBase
//	    DataID     int
//	    prevDataID int
//	    Data       *DataModel
//	}
//
//	func (c *DataDisplay) OnParametersSet() {
//	    // Manual change detection
//	    if c.DataID != c.prevDataID {
//	        c.prevDataID = c.DataID
//	        go c.fetchData()
//	    }
//	}
type ParameterReceiver interface {
	OnParametersSet()
}

// Cleaner is implemented by components that need cleanup when unmounted.
// OnDestroy is called once when the component instance is removed from the component tree.
//
// Example:
//
//	type TimerComponent struct {
//	    runtime.ComponentBase
//	    ctx    context.Context    `nojs:"state"`
//	    cancel context.CancelFunc `nojs:"state"`
//	    Count  int                `nojs:"state"`
//	}
//
//	func (c *TimerComponent) OnInit() {
//	    c.ctx, c.cancel = context.WithCancel(context.Background())
//	    go c.startTimer()
//	}
//
//	func (c *TimerComponent) startTimer() {
//	    ticker := time.NewTicker(1 * time.Second)
//	    defer ticker.Stop()
//	    for {
//	        select {
//	        case <-c.ctx.Done():
//	            return // Cleanup complete
//	        case <-ticker.C:
//	            c.Count++
//	            c.StateHasChanged()
//	        }
//	    }
//	}
//
//	func (c *TimerComponent) OnDestroy() {
//	    if c.cancel != nil {
//	        c.cancel() // Stop the timer goroutine
//	    }
//	}
type Cleaner interface {
	OnDestroy()
}

// PropUpdater is implemented by generated component code to support prop updates.
// This interface is used internally by the framework and should not be implemented manually.
// The compiler generates the ApplyProps method automatically for each component.
//
// The generated method copies props from the source component to the receiver while
// preserving internal state fields.
type PropUpdater interface {
	// ApplyProps copies props from the source component to the receiver.
	// The compiler generates this method automatically.
	ApplyProps(source Component)
}
