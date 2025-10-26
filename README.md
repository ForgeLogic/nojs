# Framework Core Philosophy

Every feature in this framework is guided by a set of core principles. These rules ensure the framework remains robust, intuitive for Go developers, and performant.

- **Type Safety Above All:** The framework will always prefer compile-time safety over runtime flexibility. The AOT compiler acts as a "Go Inspector," reading your Go code to validate props, methods, and expressions, catching errors before they ever reach the browser.
    
- **Go-Idiomatic by Default:** The template language and component architecture should feel like a natural extension of Go. Syntax and patterns (like the `for...range` directive) are modeled directly on Go's own semantics to provide a familiar and intuitive developer experience.
    
- **Explicit is Better than Implicit:** The framework favors clarity and developer control over "magic." Features like manual state updates (`StateHasChanged()`) and mandatory list keys (`trackBy`) make the flow of data and rendering predictable and easy to debug.
    
- **Simplicity Through Focus:** Features that add complexity for little practical benefit will be excluded. The goal is a lean, focused API that provides powerful tools for the most common web development patterns.
    
- **Unopinionated by Default:** The framework avoids imposing strict architectural patterns where possible. It provides core functionalities (like component rendering and lifecycle) but leaves broader concerns like state management and project structure largely to the developer, only enforcing opinions when essential for the framework's core operation (e.g., type safety checks, `trackBy` for lists).

# Project Compilation Instructions

To compile this project, please follow these steps:

1. Open a terminal and navigate to the root directory of the project.
2. Run the following command to build the project for WebAssembly:

   ``` bash
   $ GOOS=js GOARCH=wasm go build -o main.wasm
   ```

   ``` PowerShell
   PS> env:GOOS="js"; $env:GOARCH="wasm"; go build -o main.wasm
   ```

This command sets the target operating system to JavaScript (`GOOS=js`) and the architecture to WebAssembly (`GOARCH=wasm`). The output will be a `main.wasm` file, which can be used in web environments that support WebAssembly.

# Running the Project

To run the project after compilation, follow these steps:

1. In the root directory of the project, start a static file web server. For example, you can use Python's built-in HTTP server:

   ``` bash
   $ python3 -m http.server 9090
   ```

   ``` PowerShell
   PS> python -m http.server 9090
   ```

2. Open your web browser and navigate to `http://localhost:9090` to access the project.

This will serve the compiled `main.wasm` and any other static files in the project directory, allowing you to run and test the application in your browser.

> Note: In your browser's DevTools, enable "Disable cache" to force loading WebAssembly modules (e.g., main.wasm) on every refresh. For Chrome/Edge, open DevTools, go to the Network tab, and check "Disable cache" (applies while DevTools is open).

# Using the AOT Compiler (HTML Template to Go Component)

The framework includes an Ahead-of-Time (AOT) compiler for converting HTML templates into Go component code. This enables automatic generation of `Render()` methods from declarative templates.

### Workflow

1. **Create your template:**  
   Place your HTML template in the `compiler` directory. The source file must be named `input.gt.html`.

2. **Run the compiler:**  
   Use the following command from the project root to generate the Go component file:

   ```PowerShell
   PS> go run ./compiler -in compiler\input.gt.html -out ..\generated.go
   ```

   - `-in` specifies the input template file.
   - `-out` specifies the output Go file (e.g., `generated.go`).

3. **Integrate the generated component:**  
   The output file will contain a Go component with a `Render()` method based on your template. You can import and use this component in your application as usual.

> **Note:** The AOT compiler is under active development. Template syntax and features may change. See repo documentation for supported bindings and events.