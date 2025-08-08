# Copilot instructions for this repo

Purpose: Help AI agents work productively in this Go + WebAssembly (WASM) demo without guessing. Keep answers concrete, repo-specific, and runnable.

## Big picture
- This is a minimal Go → WebAssembly playground exposing Go functions to the browser via syscall/js.
- JS bootstraps the wasm (wasm_exec.js + core.js), loads HTML partials into placeholders, and calls exported Go funcs. Go can also call back into JS.
- Virtual DOM scaffold exists (`vdom/`, `component/`) and includes a tiny renderer that supports `<p>` only.

Key files
- `main.go`: wasm entrypoint. Exports `add`, calls into JS `calledFromGoWasm`, then blocks with `select {}` to keep the runtime alive.
- `core.js`: loads `main.wasm` with `Go` from `wasm_exec.js`, and provides `loadComponent(id, path)` for HTML partials.
- `labtests.js`: browser-side helpers to call exported Go (`add`) and JS callback (`calledFromGoWasm`).
- `console/`, `dialogs/`, `sessionStorage/`: thin wrappers around `syscall/js` for browser APIs (console, alert/prompt, sessionStorage).
- `vdom/vnode.go`: simple VNode type and helpers. `vdom/render.go`: minimal DOM renderer (only `<p>`). `component/component.go`: a `Component` interface with `Render() *vdom.VNode`.
- `index.html`, `header.html`, `main.html`, `footer.html`: static shell + partials.

## Build and run (local dev)
- Build wasm at repo root:
  - Env: GOOS=js GOARCH=wasm
  - Output: `main.wasm`
- Static-serve the folder (any server OK) and open `index.html`.

Example sequence
1) Build: GOOS=js GOARCH=wasm go build -o main.wasm
2) Serve: python3 -m http.server 9090
3) Browse: http://localhost:9090 (console logs show module load)

Notes
- `wasm_exec.js` is vendored (Go runtime bridge). Keep in sync with installed Go when upgrading.
- `core.js` expects `main.wasm` at project root.

## Patterns and conventions
- Build tags: All browser-facing Go files use `//go:build js || wasm` to target wasm.
- JS↔Go interop:
  - Export Go to JS by `js.Global().Set("name", js.FuncOf(fn))` (see `add` in `main.go`).
  - Call JS from Go via `js.Global().Call("fnName", args...)` (see `calledFromGoWasm`).
  - Keep Go alive with `select {}` at end of `main()`.
- Browser API wrappers: Prefer packages `console`, `dialogs`, `sessionStorage` over raw `syscall/js` in app code.
- HTML partials: Loaded into `#header`, `#content`, `#footer` with `loadComponent()`; keep paths relative to repo root when serving.
- VDOM: `vdom.VNode` has a minimal renderer; no diff/patch logic. Only `<p>` renders; others are ignored for now.

## Commit message guidelines
- When suggesting a commit message, always use the [Conventional Commits](https://www.conventionalcommits.org/) specification (e.g., `feat:`, `fix:`, `chore:`, etc.).
- The commit message must include very detailed information about the reason for the changes, not just what was changed.
- Example: `fix(component): correct VNode rendering for <p> elements to prevent double rendering in vdom/render.go. This resolves a bug where paragraphs were duplicated due to incorrect child node handling.`

## Adding features (examples)
- Expose a new Go function to JS:
  - Implement `func doThing(this js.Value, args []js.Value) interface{}` in `main.go` or a new file with wasm build tag.
  - Register with `js.Global().Set("doThing", js.FuncOf(doThing))`.
  - Call from JS: `window.doThing("arg")` after wasm has started (after core.js runs `go.run`).
- Use wrappers:
  - Logs: `console.Log("msg", 123)`; Warn/Error similar.
  - Dialogs: `dialogs.Alert("Hi")`, `name := dialogs.Prompt("Your name?")`.
  - Session storage: `sessionStorage.SetItem("k","v")`, `GetItem`, `RemoveItem`, `Clear`, `Length`.
- Compose VNodes:
  - Create nodes: `vdom.NewVNode("p", nil, nil, "hello")` or `vdom.Paragraph("hello")`
  - Mount: `vdom.RenderToSelector("#content", vdom.Paragraph("Hi"))`
  - Implement component: type MyComp struct{}; func (c MyComp) Render() *vdom.VNode { return vdom.NewVNode("div", nil, nil, "hi") }

## Debugging tips
- Open browser DevTools:
  - Console should print "WebAssembly module loaded." and logs from `main.go` calls.
  - If `add` is undefined, ensure `main.go` exported it and wasm is rebuilt/served fresh.
- Common pitfalls:
  - Not rebuilding after Go changes (always rebuild `main.wasm`).
  - Serving from wrong directory or missing `wasm_exec.js`/`core.js` includes in `index.html`.
  - Calling JS before `go.run(...)` completes; wait until the wasm runtime has started.

## Upgrades and compatibility
- Go version is declared in `go.mod`. If you upgrade Go, refresh `wasm_exec.js` to the matching version (from your Go toolchain).
- Keep build tags consistent across browser-targeted packages.

## Quick references
- Entrypoint: `main.go`
- Interop: syscall/js + wrappers in `console/`, `dialogs/`, `sessionStorage/`
- UI: `index.html` (+ partials), `core.js`, `labtests.js`
- VDOM types: `vdom/vnode.go`, API contract: `component/component.go`
