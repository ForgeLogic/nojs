# Changelog

All notable changes to the **nojs** framework are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Project Status:** nojs is currently in **MVP (Minimum Viable Product) / Experimental** phase. The API is subject to change. Use for experimentation, prototyping, and community feedback.

---

## [0.1.0-alpha] — 2026-02-23

### 🎉 Initial MVP Release

This is the first experimental release of **nojs**, a type-safe web framework for building high-performance front-end applications entirely in Go with WebAssembly and AOT-compiled HTML templates.

#### Added

##### Core Framework (`nojs/`)
- **Component Model**: Go structs combined with HTML templates for building reusable UI components
- **Component Lifecycle**: Full lifecycle support with `OnMount()`, `OnParametersSet()`, and `OnUnmount()` hooks
- **State Management**: Manual state updates via `StateHasChanged()` for predictable, explicit control flow
- **Virtual DOM (VDOM)**: In-memory DOM representation with efficient diff/patch cycles for minimal real DOM operations
- **Instance Caching**: Child components are automatically cached and reused across render cycles, preserving internal state
- **Dev vs Prod Modes**: Separate panic handling strategies via build tags for better debugging and production stability

##### AOT Template Compiler (`compiler/`)
- **HTML Template Compilation**: Parse `.gt.html` template files and auto-generate type-safe `Render()` methods at build time
- **Type-Safe Props & Methods**: Compile-time validation of component properties, methods, and event handlers
- **Data Binding**: Bind component data to DOM elements with `{PropertyName}` syntax
- **Event Binding**: Type-safe event handling with `@event="MethodName"` syntax
  - Supported events: `click`, `change`, `input`, `submit`, `focus`, `blur`, and more
- **Conditional Rendering**: `{@if}`, `{@elseif}` and `{@else}` blocks for dynamic UI based on component state
- **List Rendering**: `{@for}` loop syntax with automatic component instance management
- **Form Handling**: Native form binding with proper change/input detection
- **Inline Conditionals**: Ternary expressions in templates (e.g., `{condition ? valueA : valueB}`)
- **Multiline Text**: Support for plain text content with line break preservation

##### SPA Router (`router/`)
- **Client-Side Routing**: Full SPA (Single Page Application) routing without page reloads
- **Route Registration**: Define routes with nested layouts and page components
- **Layout Pivot Pattern**: Automatically reuse layout instances across navigations to preserve layout state
- **Route Parameters**: Extract and use route parameters (e.g., `/user/{id}`)
- **Programmatic Navigation**: `Navigate(path)` method for dynamic routing from any component
- **Browser History Integration**: Automatic `popstate` listener for browser back/forward button support
- **Nested Routes**: Support for multi-level layout hierarchies with automatic pivot reuse

##### Runtime & Core Features (`nojs/runtime`, `nojs/vdom`)
- **Virtual DOM Rendering**: Minimal DOM renderer with support for a *limited amount* of standard HTML elements
- **Component Renderer Interface**: Extensible renderer pattern for flexible DOM manipulation
- **Browser API Wrappers**:
  - `console/`: Type-safe logging (`Log`, `Warn`, `Error`)
  - `dialogs/`: Alert, confirm, and prompt dialogs
  - `sessionStorage/`: Session storage operations (`SetItem`, `GetItem`, `RemoveItem`, `Clear`, `Length`)
- **Event Registry**: Type-safe event handling system with adapter pattern for browser events

##### Development & Build Tools
- **Makefile Build System**: Unified build commands for template compilation and WASM generation
  - `make full`: Development build with all features
  - `make full-prod`: Production build with optimizations
  - `make wasm`: Fast WASM-only rebuild (~1-2 seconds)
  - `make serve`: Built-in development server (port 9090)
- **Workspace Structure**: Go workspace configuration enabling development across multiple modules (`nojs`, `compiler`, `router`, `app`)

##### Demo Application
- **Feature Showcase**: Example app demonstrating all framework capabilities
- **Example Components**:
  - Counter with state management
  - Forms with data binding
  - List rendering with TrackBy
  - Conditional rendering examples
  - Nested routing with layouts

##### Documentation
- **Framework Philosophy**: [Manifesto](https://forgelogic.github.io/nojs/design/manifesto/) explains core principles and design philosophy
- **Quick Start Guide**: [Quick Guide](https://forgelogic.github.io/nojs/guides/quick-guide/) covers all implemented features
- **Installation Guide**: [INSTALLATION.md](INSTALLATION.md) for scaffolding new projects
- **Design Decisions**: [Design Decisions](https://forgelogic.github.io/nojs/design/design-decisions/) documents key architectural choices
- **Feature Documentation**:
  - List rendering and TrackBy behavior
  - Inline conditionals support
  - Router architecture and pivot pattern
  - Text node rendering
  - Router layout engine

#### Known Limitations & Not Implemented

The following features are **not** included in this MVP release and are planned for future versions:

- **Expression Evaluation**: Only simple property binding (`{Title}`) and ternary conditionals (`{condition ? a : b}`) are supported. Complex expressions (method calls, operators, etc.) are not yet supported.
- **Advanced Router Features**: Query parameters, hash-based routing, and lazy component loading are not yet implemented.
- **CSS-in-JS**: No scoped CSS or CSS-in-Go support. Use external stylesheets.
- **Form Validation**: Built-in form validation is not included; validation must be implemented in component logic.
- **Standard Component Library**: No built-in UI library. Basic HTML elements only.
- **Hot Reloading**: Development server does not support hot module reloading.
- **Server-Side Rendering**: WASM-only. SSR is not supported.
- **Testing Utilities**: Minimal testing support in MVP. Basic test helpers available in `compiler/testcomponents/`.

#### Breaking Changes

N/A — This is the initial release.

#### Deprecated

N/A — Nothing deprecated yet.

#### Security

- Project follows standard Go security practices
- No telemetry or external tracking
- All code is auditable and open-source under Apache 2.0

#### Performance Improvements

- **AOT Compilation**: All template parsing happens at build time, not in the browser
- **Virtual DOM Diffing**: Minimal DOM mutations via efficient VDOM reconciliation
- **WebAssembly**: Near-native performance via Go-compiled WASM binaries
- **Component Caching**: Reused component instances across render cycles reduce garbage collection pressure

---

## Contributing

We welcome community feedback and contributions. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the [Apache License 2.0](LICENSE).

---

**Thank you for trying nojs!** We're excited to hear your feedback as we build toward a stable, production-ready framework. Please open issues on GitHub to report bugs, suggest features, or share your experience.
