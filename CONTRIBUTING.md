# Contributing to nojs

First off, thank you for considering contributing to **nojs**! It's people like you who will help make Go a first-class citizen for modern web development.

> **Note:** We prioritize a low-pressure, high-quality environment. This is a community project built in our collective free time.

---

## 📖 Table of Contents

1. [🕒 Your Time is Valued](#-your-time-is-valued)
2. [🎯 Finding a Task](#-finding-a-task)
3. [🏗️ Our Engineering Philosophy](#️-our-engineering-philosophy)
4. [🛠️ How Can I Contribute?](#️-how-can-i-contribute)
5. [📋 Development Workflow](#-development-workflow)
6. [⚖️ Legal & CLA](#️-legal--cla)

---

## 🕒 Your Time is Valued

We don't use strict deadlines or high-pressure roadmaps. We understand that contributors have jobs, families, and lives.

* **Work at your own pace:** There is no rush.
* **Life happens:** If you start an issue but can't finish it, just leave a quick comment. No judgment.
* **Quality over Speed:** We prefer a PR that takes three weeks to be completed over one that takes three hours and breaks our type-safety goals.

## 🎯 Finding a Task

Not sure where to start? Check our **[GitHub Project Board](https://www.google.com/search?q=https://github.com/orgs/ForgeLogic/projects/YOUR_PROJECT_ID)**.

* **`good first issue`**: Small tasks perfect for getting familiar with the AOT compiler and the runtime engine.
* **`help wanted`**: Well-defined features ready for implementation.
* **`area: compiler` or `area: runtime`**: If you have specific expertise, filter by these labels.

---

## 🏗️ Our Engineering Philosophy

To keep **nojs** lean and performant, all contributions must align with these four pillars:

### 1. Type Safety Above All

* Avoid `interface{}` (or `any`) unless absolutely necessary.
* We prefer compile-time errors over runtime flexibility.
* The AOT compiler validates templates at build time.

### 2. Go-Idiomatic Code

* Follow [Effective Go](https://go.dev/doc/effective_go).
* If a pattern doesn't feel like standard Go, it likely doesn't belong here.
* Use standard patterns: `if err != nil`, `for...range`, etc.

### 3. Minimal Dependencies

* Core remains lean. Use the Go standard library whenever possible.
* Current blessed dependencies: `golang.org/x/net`, `golang.org/x/tools`.

### 4. Explicit Over Implicit

* Manual state updates via `StateHasChanged()`.
* No hidden "magic" – predictable, debuggable data flow.

---

## 🛠️ How Can I Contribute?

### Reporting Bugs

* Check existing [Issues](https://github.com/ForgeLogic/nojs/issues).
* Provide a **minimal reproducible example** (Go code + `.gt.html` template).

### Suggesting Enhancements

* Open an issue labeled `enhancement` to discuss the idea **before** coding.
* We are selective to avoid feature creep.

### Pull Requests

1. **Fork** and branch from `main`.
2. **Follow Standards:** Run `gofmt` and use build tags (`//go:build js && wasm`).
3. **Commit Messages:** Use [Conventional Commits](https://www.conventionalcommits.org/) (e.g., `feat(compiler): add support for slots`).
4. **Testing:** - Ensure the AOT compiler still generates valid Go code for `nojs/testcomponents/`.
* Verify changes in at least one modern browser (Chrome, Firefox, or Safari). Mention which one you used in your PR description.



---

## 📋 Development Workflow

> [!IMPORTANT]
> **nojs** relies on an AOT compiler. If you change a `.gt.html` template, you **must** re-run the compiler before building the WASM binary.

```bash
# 1. Build the AOT compiler
cd nojs && go build ./cmd/nojs-compiler

# 2. Compile templates (run from your app directory)
# This generates the .go files from your .gt.html files
go run github.com/ForgeLogic/nojs/cmd/nojs-compiler -in ./components -dev

# 3. Build WASM
GOOS=js GOARCH=wasm go build -o wwwroot/main.wasm

# 4. Fast Track (using Makefile)
make full        # Compile + Build WASM (dev mode)
make serve       # Start local dev server

```

---

## ⚖️ Legal & CLA

By contributing to **nojs**, you agree that your work is original and licensed under the [Apache License 2.0](https://www.google.com/search?q=LICENSE). You retain your copyright, but grant the community an irrevocable license to use and distribute your contribution.

---

## 🤝 Code of Conduct

Be respectful. Assume good intent. Focus on constructive technical feedback. We’re all here to build something cool.

**Questions?** Start a [Discussion](https://github.com/ForgeLogic/nojs/discussions) or check the docs at [forgelogic.github.io/nojs](https://forgelogic.github.io/nojs/).
