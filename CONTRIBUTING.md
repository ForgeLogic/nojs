# Contributing to nojs

First off, thank you for considering contributing to **nojs**! It's people like you who will help make Go a first-class citizen for modern web development.

By contributing to this project, you agree to abide by its terms and the [Apache License 2.0](LICENSE).

---

## 🏗️ Our Engineering Philosophy

Before you submit a Pull Request, please ensure your contribution aligns with the **nojs** framework philosophy:

### 1. **Type Safety Above All**
- Avoid `interface{}` (or `any`) unless absolutely necessary
- We prefer compile-time errors over runtime flexibility
- Components, props, and event handlers must be strongly typed
- The AOT compiler validates templates at build time, not runtime

### 2. **Go-Idiomatic Code**
- Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- If a pattern doesn't feel like standard Go, it probably doesn't belong in **nojs**
- Use Go's standard patterns: `for...range`, `if err != nil`, struct embedding, etc.

### 3. **Minimal Dependencies**
- We aim to keep the framework core lean
- Before adding a new third-party dependency, consider if the same functionality can be achieved using the Go standard library
- Current dependencies: `golang.org/x/net`, `golang.org/x/tools`

### 4. **Explicit Over Implicit**
- Manual state updates via `StateHasChanged()`
- Mandatory `trackBy` keys in list rendering
- No hidden "magic" – predictable, debuggable data flow
- Clear component lifecycle and update triggers

---

## 🛠️ How Can I Contribute?

### Reporting Bugs

- Check the [Issues](https://github.com/ForgeLogic/nojs/issues) to see if the bug has already been reported
- If not, open a new issue with:
  - **Minimal reproducible example** (Go code + `.gt.html` template if applicable)
  - Your environment: Go version, Browser, OS
  - Expected vs. actual behavior
  - Console errors (if browser-related)

### Suggesting Enhancements

- We love new ideas, but we are highly selective to avoid "feature creep"
- Open an issue labeled `enhancement` to discuss the idea **before** spending time on implementation
- Explain:
  - How it aligns with the framework's philosophy (type safety, Go idioms, explicitness)
  - Why existing patterns can't solve the problem
  - Expected API surface and usage examples

### Pull Requests

1. **Fork** the repository and create your branch from `main`

2. **Follow Code Standards:**
   - Run `gofmt` on all Go files
   - Use build tags: `//go:build js && wasm` for browser-targeted code
   - Write tests for new functionality (see `nojs/testcomponents/` for examples)
   - Ensure all tests pass: `go test ./...`

3. **Commit Message Format:**
   - Use [Conventional Commits](https://www.conventionalcommits.org/) specification
   - Format: `type(scope): detailed description of changes with rationale`
   - Examples:
     ```
     feat(compiler): add support for nested component slots
     
     This enables parent-child content projection by detecting []*vdom.VNode 
     fields in layout components and auto-generating injection code in Render() 
     methods. Resolves #42.
     
     fix(vdom): prevent duplicate event listener registration in patch cycle
     
     Event listeners were being re-attached on every update due to missing 
     comparison logic in the differ. Now tracks listener identity to avoid 
     duplicates, reducing memory leaks in long-lived components.
     
     docs(readme): clarify AOT compiler usage with directory scanning
     
     Updated CLI examples to show that -in accepts directories, not individual 
     files. Removed incorrect -out flag documentation.
     ```

4. **Testing Requirements:**
   - If modifying the AOT compiler, ensure it generates valid Go code for existing templates in `nojs/testcomponents/`
   - Add test cases for new features (`.gt.html` templates + expected generated code)
   - Verify browser compatibility (Chrome, Firefox, Safari) for WASM changes
   - Test in both development (`-dev` flag) and production mode

5. **Update Documentation:**
   - Update `README.md` for user-facing changes
   - Update `/nojs/documentation/` files for architectural changes
   - Add inline comments for complex logic, especially in compiler and VDOM code

---

## 📋 Development Workflow

### Building and Testing Locally

```bash
# Build the AOT compiler
cd nojs
go build ./cmd/nojs-compiler

# Compile templates (from app directory)
cd ../app
go run github.com/vcrobe/nojs/cmd/nojs-compiler -in ./components -dev

# Build WASM
GOOS=js GOARCH=wasm go build -o wwwroot/main.wasm

# Or use Makefile shortcuts
make full        # Compile templates + build WASM (dev mode)
make full-prod   # Compile templates + build WASM (optimized)
make wasm        # Rebuild WASM only (fast iteration)
make serve       # Start development server
```

### Testing in Browser

1. Start server: `make serve` (or `cd app/wwwroot && python3 -m http.server 9090`)
2. Open `http://localhost:9090`
3. Check browser console for:
   - "WebAssembly module loaded." message
   - Component render logs (in `-dev` mode)
   - Any runtime errors

### Code Review Checklist

Before submitting, verify:
- [ ] Code follows Go idioms and formatting (`gofmt`)
- [ ] Build tags are correct for WASM-targeted files (`//go:build js && wasm`)
- [ ] Tests pass (`go test ./...`)
- [ ] Commit messages follow Conventional Commits format with detailed rationale
- [ ] Documentation is updated (README, `/nojs/documentation/`, inline comments)
- [ ] No new third-party dependencies (or justified in PR description)
- [ ] AOT compiler generates valid Go code for test templates
- [ ] Browser testing completed (Chrome, Firefox, Safari)

---

## ⚖️ Contributor License Agreement (CLA)

By submitting a Pull Request to **nojs**, you agree that:

- Your contribution is your original work
- You grant **ForgeLogic** a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright and patent license to use and distribute your contribution as part of the project
- Your contribution is licensed under the [Apache License 2.0](LICENSE)

This is a standard requirement for open-source projects under Apache 2.0 – you retain copyright, but grant usage rights to the community.

---

## 🤝 Code of Conduct

We are committed to providing a welcoming and inclusive environment:

- Be respectful and professional in all interactions
- Focus on constructive feedback
- Assume good intent
- Harassment and abusive behavior will not be tolerated

Report violations to the project maintainers.

---

## 💡 Questions?

- Open a [Discussion](https://github.com/ForgeLogic/nojs/discussions) for general questions
- Use [Issues](https://github.com/ForgeLogic/nojs/issues) for bugs and feature requests
- Check existing documentation in `/nojs/documentation/`:
  - [`LIST_RENDERING.md`](nojs/documentation/LIST_RENDERING.md) - List optimization with `trackBy`
  - [`INLINE_CONDITIONALS.md`](nojs/documentation/INLINE_CONDITIONALS.md) - Template conditional logic
  - [`ROUTER_ARCHITECTURE.md`](nojs/documentation/ROUTER_ARCHITECTURE.md) - SPA routing system

Thank you for contributing to **nojs**! 🚀
