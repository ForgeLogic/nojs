package appstate

import (
	"github.com/ForgeLogic/nojs/signals"
)

// App-level signals — owned by the app, not the framework.
// Add new global signals here as the app grows.
var RenderCount = signals.NewSignal(1)
var NextIDIndex = signals.NewSignal(1)
