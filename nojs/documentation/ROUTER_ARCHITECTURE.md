# Router Architecture and Implementation

## Overview

The No-JS framework implements a **pluggable, router-agnostic architecture** for client-side routing. The router integrates seamlessly with the Virtual DOM (VDOM) and component lifecycle system to enable Single Page Application (SPA) navigation without full page reloads.

This document covers all technical details of the router implementation, including architecture decisions, integration patterns, event handling, and VDOM patching challenges.

---

## Table of Contents

1. [Architecture and Design Principles](#architecture-and-design-principles)
2. [Core Interfaces](#core-interfaces)
3. [Router Implementation](#router-implementation)
4. [Routing Modes](#routing-modes)
5. [Route Matching and Parameter Extraction](#route-matching-and-parameter-extraction)
6. [Integration with Renderer](#integration-with-renderer)
7. [Component Navigation](#component-navigation)
8. [Event System Integration](#event-system-integration)
9. [VDOM Event Listener Management](#vdom-event-listener-management)
10. [Browser History Integration](#browser-history-integration)
11. [Lifecycle and Initialization](#lifecycle-and-initialization)
12. [Usage Examples](#usage-examples)
13. [Technical Challenges and Solutions](#technical-challenges-and-solutions)

---

## Architecture and Design Principles

### Design Philosophy

The router architecture follows three key principles:

1. **Router Agnostic**: The framework core (`runtime` package) doesn't depend on any specific router implementation
2. **Pluggable**: Any router that implements the `NavigationManager` interface can be used
3. **Unopinionated**: Supports multiple routing strategies (path-based, hash-based) without prescribing one approach

### Separation of Concerns

```
┌─────────────────┐
│   Application   │
│    (main.go)    │
└────────┬────────┘
         │
         ├──────────────────┐
         │                  │
         ▼                  ▼
┌────────────────┐   ┌──────────────┐
│    Router      │   │   Renderer   │
│   (router/)    │◄──┤  (runtime/)  │
└────────────────┘   └──────┬───────┘
         │                   │
         │                   ▼
         │            ┌──────────────┐
         │            │  Components  │
         │            │ (ComponentBase)│
         │            └──────────────┘
         │
         ▼
┌────────────────┐
│  Browser APIs  │
│ (History/Hash) │
└────────────────┘
```

The router is **injected** into the renderer at initialization, allowing the framework to work with or without routing.

---

## Core Interfaces

### NavigationManager Interface

Defined in `runtime/navigation.go`, this is the contract that any router must implement:

```go
type NavigationManager interface {
    // Start initializes the router with an onChange callback
    Start(onChange func(chain []Component, key string)) error
    
    // Navigate programmatically changes the URL and renders new component
    Navigate(path string) error
    
    // GetComponentForPath resolves a path to its component
    GetComponentForPath(path string) (Component, bool)
}
```

**Key Responsibilities**:
- Initialize browser event listeners (popstate for back/forward buttons)
- Read initial URL on application startup
- Match URL paths to registered routes
- Create component instances via route handlers
- Call the `onChange` callback to trigger rendering with a chain of components and a unique key

### Navigator Interface

Defined in `runtime/navigation.go`, this is provided to components:

```go
type Navigator interface {
    Navigate(path string) error
}
```

**Implementation Chain**:
```
Component.Navigate() → ComponentBase.Navigate() → Renderer.Navigate() → Router.Navigate()
```

This chain allows components to trigger navigation without direct coupling to the router.

---

## Router Implementation

### Router Structure

Located in `router/router.go`:

```go
type Router struct {
    routes           []routeDefinition
    onChange         func(runtime.Component, string) // Second parameter is path/key
    mode             RoutingMode
    notFoundHandler  RouteHandler
    popstateListener js.Func
}

type routeDefinition struct {
    Path    string       // e.g., "/users/{id}"
    Handler RouteHandler // func(params map[string]string) runtime.Component
}

type RouteHandler func(params map[string]string) runtime.Component
```

**Design Notes**:
- `RouteHandler` is a factory function that creates component instances
- Handlers receive extracted URL parameters (e.g., `{id}` becomes `params["id"]`)
- `popstateListener` is a `js.Func` that must be properly released to avoid memory leaks

### Route Registration

Routes are registered via `Handle()`:

```go
appRouter.Handle("/users/{id}", func(params map[string]string) runtime.Component {
    return &UserPage{UserID: params["id"]}
})
```

**Important**: Route handlers are called on **every navigation** to create fresh component instances. This ensures clean state for each route. Route handlers receive extracted URL parameters as a map.

### 404 Handling

Optional not-found handler:

```go
appRouter.HandleNotFound(func(params map[string]string) runtime.Component {
    return &NotFoundPage{}
})
```

If no handler is registered and no route matches, the router prints a warning but doesn't crash.

---

## Routing Modes

### PathMode (Default)

Uses the **HTML5 History API** with clean URLs:

```
https://example.com/about
https://example.com/users/123
```

**Browser API Used**: `history.pushState()`

**Server Requirement**: Server must be configured to serve `index.html` for all routes (SPA fallback).

**Example Server Configs**:
- **Go http.FileServer**: Use `http.StripPrefix` with custom fallback handler
- **Nginx**: `try_files $uri /index.html`
- **Apache**: `RewriteRule` to `index.html`

### HashMode

Uses hash-based URLs (no server config needed):

```
https://example.com/#/about
https://example.com/#/users/123
```

**Browser API Used**: `location.hash`

**Advantages**:
- Works with static file hosting (GitHub Pages, S3, etc.)
- No server configuration required
- Browser natively handles hash navigation

**Implementation**:

```go
appRouter := router.New(&router.Config{Mode: router.HashMode})
```

The router automatically strips the `#` prefix when parsing paths.

---

## Route Matching and Parameter Extraction

### Pattern Matching

The `matchRoute()` function handles pattern matching with parameter extraction:

```go
func (r *Router) matchRoute(routePath, urlPath string) (map[string]string, bool)
```

**Algorithm**:

1. **Normalize paths**: Remove trailing slashes, handle empty strings as `/`
2. **Split into segments**: Split on `/` delimiter
3. **Length check**: Routes must have same number of segments
4. **Segment-by-segment comparison**:
   - Static segments must match exactly
   - Dynamic segments (wrapped in `{}`) capture the URL value
5. **Return** extracted parameters and match status

**Examples**:

```go
// Static route
matchRoute("/about", "/about") 
→ (map[]{}, true)

// Dynamic route
matchRoute("/users/{id}", "/users/123") 
→ (map["id": "123"], true)

// Multi-parameter route
matchRoute("/posts/{year}/{month}/{slug}", "/posts/2024/11/hello") 
→ (map["year": "2024", "month": "11", "slug": "hello"], true)

// No match
matchRoute("/about", "/contact") 
→ (nil, false)
```

### Parameter Constraints

**Currently supported**:
- Dynamic segments with curly braces (e.g., `{id}`, `{year}`, `{slug}`)
- Multiple parameters per route (e.g., `/posts/{year}/{month}/{slug}`)
- Parameter extraction passed to component factories

**NOT yet supported** (future enhancement):
- Type constraints (e.g., `{id:int}`)
- Regex patterns (e.g., `{slug:[a-z-]+}`)
- Optional segments
- Wildcard routes

---

## Integration with Renderer

### Renderer Initialization

The renderer accepts an optional `NavigationManager`:

```go
renderer := runtime.NewRenderer(appRouter, "#app")
```

If `nil` is passed, the renderer works without routing (useful for non-SPA apps or embedded components).

### onChange Callback

The application defines how to respond to navigation:

```go
onRouteChange := func(newComponent runtime.Component, path string) {
    renderer.SetCurrentComponent(newComponent, path)
    renderer.ReRender()
}

appRouter.Start(onRouteChange)
```

**Execution Flow**:

```
URL Change → Router.handlePathChange() 
          → GetComponentForPath() 
          → RouteHandler(params) 
          → onChange(newComponent, path) 
          → Renderer.SetCurrentComponent() 
          → Renderer.ReRender()
          → VDOM Patching
```

### SetCurrentComponent

Located in `runtime/renderer_impl.go`:

```go
func (r *RendererImpl) SetCurrentComponent(comp Component, key string) {
    r.currentComponent = comp
    r.currentKey = key
}
```

This swaps out the root component **without destroying the renderer instance**, preserving:
- Component instance cache (`r.instances`)
- Previous VDOM tree (`r.prevVDOM`)
- Lifecycle tracking (`r.initialized`, `r.activeKeys`)

The `key` parameter helps the renderer track component identity for efficient reconciliation.

---

## Component Navigation

### ComponentBase.Navigate()

Every component that embeds `runtime.ComponentBase` can trigger navigation:

```go
type MyComponent struct {
    runtime.ComponentBase
}

func (c *MyComponent) HandleClick() {
    c.Navigate("/about")
}
```

**Implementation** (`runtime/componentbase.go`):

```go
func (b *ComponentBase) Navigate(path string) error {
    if b.renderer == nil {
        return fmt.Errorf("renderer is nil (component not mounted?)")
    }
    return b.renderer.Navigate(path)
}
```

**Flow**:

1. Component calls `Navigate(path)`
2. ComponentBase delegates to `renderer.Navigate(path)`
3. Renderer delegates to `navManager.Navigate(path)`
4. Router updates browser URL and calls `onChange` callback
5. Renderer re-renders with new component

**Error Handling**: Returns error if renderer not set (component not mounted yet).

---

## Event System Integration

### EventBase Composition Pattern

All event argument types embed `EventBase` to provide common functionality:

```go
type EventBase struct {
    jsEvent               js.Value
    preventDefaultCalled  bool
    stopPropagationCalled bool
}

func (e *EventBase) PreventDefault() {
    if !e.preventDefaultCalled {
        e.jsEvent.Call("preventDefault")
        e.preventDefaultCalled = true
    }
}
```

### ClickEventArgs

Used by Link component and other click handlers:

```go
type ClickEventArgs struct {
    EventBase  // Embedded for PreventDefault/StopPropagation
    Button     int
    ClientX    int
    ClientY    int
    AltKey     bool
    CtrlKey    bool
    ShiftKey   bool
}
```

### Dual Signature Support

The `onclick` event supports two handler signatures:

```go
// No arguments (for simple actions)
func (c *Component) HandleClick() { ... }

// With event args (for advanced handling)
func (c *Component) HandleClick(e events.ClickEventArgs) { ... }
```

This is validated at **compile time** by the AOT compiler (`compiler/compiler.go`).

### Event Adapters

Located in `events/adapters.go`, these convert Go functions to JavaScript callbacks:

```go
func AdaptClickEvent(handler func(events.ClickEventArgs)) func(js.Value) {
    return func(jsEvent js.Value) {
        eventBase := NewEventBase(jsEvent)
        args := ClickEventArgs{
            EventBase: eventBase,
            Button:    jsEvent.Get("button").Int(),
            ClientX:   jsEvent.Get("clientX").Int(),
            ClientY:   jsEvent.Get("clientY").Int(),
            // ... extract other properties
        }
        handler(args)
    }
}
```

**Flow**:

```
DOM Event → js.FuncOf wrapper → AdaptClickEvent 
         → Extract event properties 
         → Call Go handler with ClickEventArgs
```

---

## VDOM Event Listener Management

### The Challenge

One of the most complex aspects of the router implementation was handling event listeners during VDOM patching. The problem:

**Naive approach**: Just call `addEventListener` during patching
**Result**: Event listeners accumulate on every navigation, causing handlers to fire multiple times

### Root Cause

JavaScript's `addEventListener()` **does not remove old listeners automatically**. When the same element is patched multiple times:

```javascript
element.addEventListener('click', handler1);  // Navigation 1
element.addEventListener('click', handler2);  // Navigation 2
// Now clicking fires BOTH handler1 and handler2!
```

### Solutions Considered

#### ❌ Solution 1: Track and Remove Individual Listeners

```go
// Store references to js.Func callbacks
// Call removeEventListener for each old listener
// Add new listeners
```

**Problems**:
- Must maintain a separate map of element → listeners
- `js.Func` references must be stored to call `.Release()`
- Complex bookkeeping, error-prone

#### ❌ Solution 2: Compare Old and New Handlers

```go
if oldVNode.Attributes["onclick"] != newVNode.Attributes["onclick"] {
    // Only re-attach if changed
}
```

**Problems**:
- **Functions cannot be compared in Go** (`panic: comparing uncomparable type func(js.Value)`)
- Even with workarounds, determining "sameness" is impossible (closures have different addresses)

#### ✅ Solution 3: Clone Element to Remove All Listeners

**Implementation** (`vdom/render.go`):

```go
func patchElement(domElement js.Value, oldVNode, newVNode *VNode) {
    // ... update attributes first ...
    
    // Check if new VNode has event handlers
    hasEventHandlers := false
    if newVNode.Attributes != nil {
        for key := range newVNode.Attributes {
            if len(key) > 2 && key[0] == 'o' && key[1] == 'n' {
                hasEventHandlers = true
                break
            }
        }
    }
    
    // Clone element to remove all listeners
    if hasEventHandlers {
        cloned := domElement.Call("cloneNode", false)  // false = don't clone children
        
        // Move children to cloned element
        for domElement.Get("firstChild").Truthy() {
            cloned.Call("appendChild", domElement.Get("firstChild"))
        }
        
        // Replace in DOM
        parent := domElement.Get("parentNode")
        if parent.Truthy() {
            parent.Call("replaceChild", cloned, domElement)
        }
        
        // Attach fresh listeners to cloned element
        attachEventListeners(cloned, newVNode.Attributes)
        
        return  // Skip remaining patching since children already moved
    }
    
    // ... continue with normal patching ...
}
```

**Why This Works**:

1. `cloneNode(false)` creates a **shallow clone** without children or event listeners
2. We manually move children from original to clone using `appendChild`
3. `replaceChild()` swaps the elements in the DOM
4. We attach **fresh listeners** to the clean clone
5. The original element (with accumulated listeners) is garbage collected

**Performance Note**: Cloning is surprisingly efficient in modern browsers. The overhead is minimal compared to the cost of event handler bugs.

### attachEventListeners Implementation

Located in `vdom/render.go`:

```go
func attachEventListeners(domElement js.Value, attributes map[string]any) {
    if attributes == nil {
        return
    }

    for key, value := range attributes {
        if len(key) > 2 && key[0] == 'o' && key[1] == 'n' {
            eventType := key[2:]  // "onclick" → "click"
            
            // Convert Go handler to JavaScript callback
            handler, ok := value.(func(js.Value))
            if !ok {
                continue
            }
            
            cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
                if len(args) > 0 {
                    handler(args[0])
                }
                return nil
            })
            
            domElement.Call("addEventListener", eventType, cb)
            
            // TODO: Store cb somewhere to release later if needed
        }
    }
}
```

**Note**: The `js.FuncOf` callbacks are currently **not explicitly released**. This is acceptable because:
- They live as long as the DOM element exists
- When the element is removed from DOM, it becomes unreachable
- Go's garbage collector will eventually clean them up
- For long-running SPAs, a future enhancement could track and release them

---

## Browser History Integration

### Popstate Event Listener

The router listens for browser back/forward button clicks:

```go
func (r *Router) Start(onChange func(runtime.Component, string)) error {
    r.onChange = onChange

    // Register popstate listener
    r.popstateListener = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        r.handlePathChange()
        return nil
    })
    js.Global().Set("onpopstate", r.popstateListener)

    // Handle initial page load
    r.handlePathChange()
    return nil
}
```

**Flow**:

```
User clicks back button → Browser fires 'popstate' event 
                       → Router.handlePathChange() 
                       → GetComponentForPath(currentURL) 
                       → onChange(newComponent, path) 
                       → Renderer re-renders
```

### Cleanup

The router provides cleanup to release the listener:

```go
func (r *Router) Cleanup() {
    if !r.popstateListener.IsUndefined() {
        r.popstateListener.Release()
    }
}
```

**Important**: In typical WASM applications that run for the entire page lifetime, cleanup is rarely needed. However, it's essential for:
- Testing scenarios
- Hot-reloading during development
- Embedding WASM modules that can be unloaded

---

## Lifecycle and Initialization

### Application Startup Sequence

1. **main.go**: Create router with config
2. **main.go**: Register routes via `Handle()`
3. **main.go**: Create renderer with router
4. **main.go**: Define `onRouteChange` callback
5. **main.go**: Call `appRouter.Start(onRouteChange)`
6. **Router**: Read initial browser URL
7. **Router**: Call `handlePathChange()`
8. **Router**: Match route and create component
9. **Router**: Call `onChange(component, path)`
10. **Renderer**: Call `SetCurrentComponent(component, path)`
11. **Renderer**: Call `ReRender()`
12. **Renderer**: Inject renderer reference into component via `SetRenderer()`
13. **Renderer**: Call component lifecycle methods (`OnInit`, `OnParametersSet`)
14. **Renderer**: Call `component.Render()`
15. **VDOM**: Render initial DOM
16. **Application**: Enter event loop (`select {}`)

### Navigation Sequence

1. **User action**: Call `component.Navigate()` from an event handler
2. **ComponentBase.Navigate()**: Delegate to `renderer.Navigate()`
3. **Renderer.Navigate()**: Delegate to `router.Navigate()`
4. **Router.Navigate()**: Call `history.pushState()` or set `location.hash`
5. **Router.Navigate()**: Call `handlePathChange()`
6. **Router.handlePathChange()**: Match path and create component
8. **Router.handlePathChange()**: Call `onChange(newComponent, path)`
9. **Renderer**: Call `SetCurrentComponent(component, path)` and `ReRender()`
9. **Renderer**: Call component lifecycle methods
10. **VDOM**: Patch DOM with minimal changes
11. **VDOM**: Clone elements with event handlers
12. **VDOM**: Attach fresh event listeners

---

## Usage Examples

### Basic Setup

```go
func main() {
    // Create router
    appRouter := router.New(&router.Config{Mode: router.PathMode})
    
    // Register routes
    appRouter.Handle("/", func(params map[string]string) runtime.Component {
        return &HomePage{}
    })
    
    appRouter.Handle("/about", func(params map[string]string) runtime.Component {
        return &AboutPage{}
    })
    
    // Create renderer
    renderer := runtime.NewRenderer(appRouter, "#app")
    
    // Start router
    appRouter.Start(func(newComponent runtime.Component, path string) {
        renderer.SetCurrentComponent(newComponent, path)
        renderer.ReRender()
    })
    
    select {}
}
```

### Parametric Routes

```go
appRouter.Handle("/users/{id}", func(params map[string]string) runtime.Component {
    userID := params["id"]
    return &UserProfilePage{UserID: userID}
})

appRouter.Handle("/blog/{year}", func(params map[string]string) runtime.Component {
    year := 2026 // Default value
    if yearStr, ok := params["year"]; ok {
        if parsed, err := strconv.Atoi(yearStr); err == nil {
            year = parsed
        }
    }
    return &BlogPage{Year: year}
})

appRouter.Handle("/posts/{year}/{month}/{slug}", func(params map[string]string) runtime.Component {
    return &BlogPostPage{
        Year:  params["year"],
        Month: params["month"],
        Slug:  params["slug"],
    }
})
```

### Component with Navigation

```go
type AboutPage struct {
    runtime.ComponentBase
}

func (a *AboutPage) NavigateToHome(e events.ClickEventArgs) {
    e.PreventDefault()
    a.Navigate("/")
}

func (a *AboutPage) Render(r *runtime.Renderer) *vdom.VNode {
    return vdom.Div(nil,
        vdom.H1(nil, "About Page"),
        vdom.A(map[string]any{
            "href": "/",
            "onclick": events.AdaptClickEvent(a.NavigateToHome),
        }, "Back to Home"),
    )
}
```

---

## Technical Challenges and Solutions

### Challenge 1: Function Comparison in Go

**Problem**: Go doesn't allow comparing functions with `==` or `!=`

**Solution**: Don't compare handlers at all. Always re-attach listeners when they exist by cloning the element.

### Challenge 2: Event Listener Accumulation

**Problem**: `addEventListener` doesn't remove old listeners

**Solution**: Clone element to strip all listeners before attaching new ones.

### Challenge 3: Preserving Component State Across Navigation

**Problem**: Creating new component instances on every navigation loses state

**Solution**: The renderer maintains an instance cache (`r.instances`) keyed by component location in the tree. Only the **root component** changes during navigation; child components are reused if they remain in the tree.

### Challenge 4: Server Configuration for PathMode

**Problem**: Direct URL access (e.g., `example.com/about`) returns 404 without server config

**Solution**: 
- Document server requirements clearly
- Provide HashMode as alternative for static hosting
- Example server configs in documentation

### Challenge 5: Preventing Memory Leaks from js.Func

**Problem**: Every `js.FuncOf` creates a callback that must be released

**Solution**: 
- **Current**: Cloning elements naturally garbage-collects old listeners
- **Future**: Implement explicit tracking and release mechanism
- **Cleanup**: Provide `Router.Cleanup()` for popstate listener

---

## Future Enhancements

### Phase 2: Advanced Route Matching

- **Type constraints**: `/users/{id:int}`
- **Regex patterns**: `/posts/{slug:[a-z0-9-]+}`
- **Optional segments**: `/search/{query?}`
- **Wildcard routes**: `/files/*filepath`
- **Query string parsing**: `/search?q=term&page=2`

### Phase 3: Navigation Guards

```go
appRouter.BeforeNavigate(func(from, to string) bool {
    if !user.IsAuthenticated() && isProtectedRoute(to) {
        return false  // Block navigation
    }
    return true
})
```

### Phase 4: Route Metadata

```go
appRouter.Handle("/admin", handler).WithMeta(map[string]any{
    "requiresAuth": true,
    "title": "Admin Panel",
})
```

### Phase 5: Lazy Loading

```go
appRouter.Handle("/admin", router.Lazy(func() runtime.Component {
    // Load admin module on demand
    return loadAdminModule()
}))
```

---

## Performance Considerations

### VDOM Patching with Event Listeners

- **Cloning overhead**: Minimal in modern browsers (~1-2ms for typical elements)
- **Trade-off**: Slight performance cost for correctness and simplicity
- **Optimization**: Only clone when `hasEventHandlers` is true

### Route Matching

- **Algorithm**: O(n) where n = number of registered routes
- **Typical usage**: Small number of routes (< 50), negligible impact
- **Future optimization**: Trie-based matching for large route tables

### Component Instance Caching

- **Strategy**: Preserve instances by key across renders
- **Benefit**: State persistence without re-initialization
- **Memory**: Instances cleaned up when unmounted via `cleanupUnmountedComponents()`

---

## Conclusion

The No-JS framework's router architecture achieves its design goals:

✅ **Pluggable**: NavigationManager interface allows custom router implementations  
✅ **Unopinionated**: Supports both PathMode and HashMode  
✅ **Integrated**: Seamless VDOM and lifecycle integration  
✅ **Correct**: Proper event listener cleanup prevents bugs  
✅ **Developer-Friendly**: Simple API with clear patterns  

The implementation handles the complexities of browser APIs, event management, and VDOM patching while exposing a clean, type-safe API to framework users.
