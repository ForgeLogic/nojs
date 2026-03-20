# Hello, App! — nojs example

The simplest possible nojs app. Type a name, get a greeting — no JavaScript written.

## What you'll learn

- How a nojs component is structured
- How to bind an event handler with `@oninput`
- How to display reactive state with `{Name}`
- How `StateHasChanged()` triggers a re-render

## Project structure

```
hello-app/
├── main.go                          # Entry point: wires up the renderer and router
├── components/
│   ├── helloapp.go                # Component struct and event handler
│   └── HelloApp.gt.html           # Template (compiled to Go by nojsc)
├── wwwroot/
│   ├── index.html                   # The single HTML page — just mounts <div id="app">
│   ├── core.js                      # Loads the WASM binary
│   └── wasm_exec.js                 # Go's WASM runtime (you copy this manually, see below)
├── go.mod
└── Makefile
```

## How to run

**1. Copy `wasm_exec.js` from your Go installation:**

```bash
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" wwwroot/
```

> On Windows with Git Bash:
> ```bash
> cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" wwwroot/
> ```

**2. Build (compile templates + WASM):**

```bash
make full
```

**3. Serve:**

```bash
make serve
```

**4. Open** http://localhost:9090

> **Windows note:** if `python3` is not in your PATH, create a `Makefile.local` file with:
> ```
> SERVE_CMD=python -m http.server 9090
> ```

---

## How it works

nojs has three pieces that work together:

### 1. The struct (`helloapp.go`)

Your component is a plain Go struct. The fields are the state.

```go
type HelloApp struct {
    runtime.ComponentBase  // required — gives the component its lifecycle methods

    Name    string
    HasName bool
}
```

### 2. The template (`HelloApp.gt.html`)

Write HTML with Go values and event bindings. The `nojsc` compiler reads this file
and generates a `Render` method automatically — you never write it by hand.

```html
<input type="text" @oninput="HandleInput" placeholder="Your name..." />

{@if HasName}
    <p>Hello, {Name}!</p>
{@else}
    <p>Hello, App!</p>
{@endif}
```

- `{Name}` — renders the value of the `Name` field
- `@oninput="HandleInput"` — calls `HandleInput` on every keystroke
- `{@if HasName} ... {@endif}` — conditional rendering

### 3. The handler

Mutate a field, then call `StateHasChanged()`. nojs diffs the new virtual DOM
against the previous one and patches only what changed — no full page reload,
no manual DOM queries.

```go
func (c *HelloApp) HandleInput(e events.ChangeEventArgs) {
    c.Name = e.Value
    c.HasName = c.Name != ""
    c.StateHasChanged()
}
```

### 4. The entry point (`main.go`)

Wires the renderer, router, and component together. For a single-page app like this
one, you register one route (`"/"`) pointing to your component.

```go
routerEngine.RegisterRoutes([]router.Route{
    {
        Path: "/",
        Chain: []router.ComponentMetadata{
            {
                Factory: func(p map[string]string) runtime.Component {
                    return &components.HelloApp{}
                },
                TypeID: helloAppTypeID,
            },
        },
    },
})
```

---

## Next steps

- Add a second component and a second route
- Try `@onclick` instead of `@oninput` (use a button to submit the name)
- Look at the other examples in `../../app` for conditionals, lists, slots, and more
