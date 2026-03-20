.PHONY: full wasm compile serve clean

COMPILER_PATH := github.com/ForgeLogic/nojs-compiler/cmd/nojsc
COMPONENTS_DIR := ./components
WASM_OUTPUT    := ./wwwroot/main.wasm

# Override in Makefile.local if needed (e.g. SERVE_CMD=python -m http.server 9090)
SERVE_CMD := python3 -m http.server 9090
SERVE_DIR := ./wwwroot

-include Makefile.local

# Full build: compile templates then WASM
full: compile wasm

# Run the nojsc compiler — generates a Render method from each .gt.html file
compile:
	@echo "Compiling templates..."
	@go run $(COMPILER_PATH) -in=$(COMPONENTS_DIR)

# Build the WASM binary
wasm:
	@echo "Building WASM..."
	@GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT) .

# Start a local HTTP server
serve:
	@cd $(SERVE_DIR) && $(SERVE_CMD)

clean:
	@rm -f $(WASM_OUTPUT)
