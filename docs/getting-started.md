# Getting Started with nojs

> ⚠️ **Status:** nojs is an experimental MVP. The API may change before a stable release.

Welcome to **nojs** — a type-safe web framework for building browser applications entirely in Go using WebAssembly.

This guide will get you up and running with a working demo in under 5 minutes. If you're ready to start your own project, head to the [Installation Guide](installation.md).

---

## What is nojs?

nojs lets you write frontend applications in **Go** that compile to WebAssembly and run in the browser at near-native speed. You get:

- **Type-safe templates** — Props, methods, and expressions are validated at compile time
- **Component model** — Reusable UI building blocks combining Go logic with HTML templates
- **Virtual DOM** — Efficient diffing and patching for fast, minimal DOM updates
- **Almost No JavaScript required** — 99% of the app code is Go

The framework follows Go's philosophy: simple, explicit, and focused. No magic. No mystery.

---

## Prerequisites

- **Go 1.25+** — [https://go.dev/dl/](https://go.dev/dl/)
- **Make** — pre-installed on Linux/macOS; on Windows use WSL
- **A static file server** — `python3 -m http.server` works fine, or any alternative

---

## 30-second quickstart

Clone and run the demo app:

```bash
git clone https://github.com/ForgeLogic/nojs.git
cd nojs
make full
make serve
```

Open **http://localhost:9090** in your browser.

### What you're seeing

The console should print `WebAssembly module loaded.` The demo app showcases:

- **Routing** — Navigate between pages without full page reloads
- **Components** — Reusable UI blocks with state and event handling
- **List rendering** — Efficient rendering of dynamic lists
- **Forms** — Two-way data binding and validation
- **Layouts** — Shared navigation and structure across pages

Expand the browser's DevTools → Console to see logs from the Go runtime.

---

## Architecture at a glance

### Components

A **component** is a reusable UI block combining three things:

1. **Go struct** — Holds state and implements event handlers
2. **HTML template** (in `*.gt.html` files) — Defines the visual structure
3. **Render() method** — Returns a virtual DOM tree (auto-generated from the template)

Example:

```go
type Counter struct {
    Count int
}

func (c *Counter) Increment() {
    c.Count++
    c.StateHasChanged() // Tell the framework to re-render
}
```

With a template (`counter.gt.html`):

```html
<div>
    <p>Count: {Count}</p>
    <button @click="Increment">Click me</button>
</div>
```

The **AOT compiler** (nojsc) reads the template and generates a type-safe `Render()` method. Props, methods, and expressions are validated at **build time**, not runtime.

### Virtual DOM

When a component's state changes, nojs:

1. Calls `Render()` to create a new virtual DOM tree
2. Diffs the new tree against the previous one
3. Calculates the minimal set of real DOM changes
4. Applies changes via WebAssembly

This minimizes expensive browser operations and keeps your app fast.

### Routing

nojs includes a built-in SPA router. Define routes in your app and navigate without full page reloads:

```go
router.Mount("/", &pages.Home{})
router.Mount("/about", &pages.About{})
```

The router supports nested routes and layout components for shared structure.

---

## Next steps

### Learn by example

Browse the demo app's components at `app/internal/app/components/` to see working patterns:

- **Pages** — Full-page components with routing
- **Dialogs** — Modal components with event handling
- **Forms** — Data binding and validation patterns
- **Lists** — Efficient rendering with `trackBy` keys

### Read the guides

- **[Installation Guide](installation.md)** — Set up your own project from scratch
- **[Quick Guide](guides/quick-guide.md)** — A practical reference for every framework feature
- **[NoJS Manifesto](design/manifesto.md)** — The design philosophy behind nojs

### Understand the framework

- **[Architecture Overview](architecture/overview.md)** — How rendering, diffing, and lifecycle work
- **[AOT Compiler Guide](compiler/how-it-works.md)** — How templates become Go code
- **[Router Internals](router/internals.md)** — How SPA navigation works

---

## Key concepts

### StateHasChanged()

After updating component state, call `StateHasChanged()` to trigger a re-render:

```go
func (c *Counter) Increment() {
    c.Count++
    c.StateHasChanged() // Framework will call Render() and patch the DOM
}
```

This is explicit and direct — you always know when the UI will update.

### Type-Safe Templates

Templates are validated at compile time. If you reference a field that doesn't exist or bind to a method that's not exported, the compiler will catch it **before deployment**. No runtime surprises.

### AOT Compilation

All template parsing happens at build time (`make full`). Nothing is parsed or compiled in the browser. This gives you predictable performance.

---

## Ready to build?

→ **[Installation Guide](installation.md)** — Start your own project

For detailed feature documentation, see the **[Quick Guide](guides/quick-guide.md)**.

For design philosophy, read the **[NoJS Manifesto](design/manifesto.md)**.