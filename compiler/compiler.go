package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/tools/go/packages"
)

// componentSchema holds the type information for a component's props.
type componentSchema struct {
	Props   map[string]propertyDescriptor // Map of Prop name to its Go type (e.g., "Title": "string")
	Methods map[string]bool               // Set of available method names for event handlers
}

type propertyDescriptor struct {
	Name          string
	LowercaseName string
	GoType        string
}

// componentInfo holds all discovered information about a component.
type componentInfo struct {
	Path          string
	PascalName    string
	LowercaseName string
	PackageName   string
	Schema        componentSchema
}

// compileOptions holds compiler-wide options passed from CLI flags.
type compileOptions struct {
	DevWarnings bool // Enable development warnings in generated code
}

// loopContext holds information about variables available in a loop scope.
type loopContext struct {
	IndexVar string // e.g., "i" or "_"
	ValueVar string // e.g., "user"
}

// Regex to find data binding expressions like {FieldName} or {user.Name}
var dataBindingRegex = regexp.MustCompile(`\{([a-zA-Z0-9_.]+)\}`)

// Regex to find ternary expressions like { condition ? 'value1' : 'value2' }
var ternaryExprRegex = regexp.MustCompile(`\{\s*(!?)([a-zA-Z0-9_]+)\s*\?\s*'([^']*)'\s*:\s*'([^']*)'\s*\}`)

// Regex to find boolean shorthand like {condition} or {!condition}
var booleanShorthandRegex = regexp.MustCompile(`^\{\s*(!?)([a-zA-Z0-9_]+)\s*\}$`)

// Standard HTML boolean attributes
var standardBooleanAttrs = map[string]bool{
	"disabled":       true,
	"checked":        true,
	"readonly":       true,
	"required":       true,
	"autofocus":      true,
	"autoplay":       true,
	"controls":       true,
	"loop":           true,
	"muted":          true,
	"selected":       true,
	"hidden":         true,
	"multiple":       true,
	"novalidate":     true,
	"open":           true,
	"reversed":       true,
	"scoped":         true,
	"seamless":       true,
	"sortable":       true,
	"truespeed":      true,
	"default":        true,
	"ismap":          true,
	"formnovalidate": true,
}

// preprocessFor preprocesses template source to extract for-loop blocks and replace them with placeholder nodes.
// It validates that every {@for} has a matching {@endfor} and that trackBy is specified.
// Syntax: {@for index, value := range SliceName trackBy uniqueKeyExpression}{@endfor}
// The index can be _ to ignore it: {@for _, value := range SliceName trackBy uniqueKeyExpression}
func preprocessFor(src string, templatePath string) (string, error) {
	// Regex to match ONLY: {@for i, user := range Users trackBy user.ID} or {@for _, user := range Users trackBy user.ID}
	reFor := regexp.MustCompile(`\{\@for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*,\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*:=\s*range\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+trackBy\s+([a-zA-Z0-9_.]+)\}`)

	// Regex to detect INVALID syntax: {@for user := range Users trackBy user.ID} (missing index/underscore)
	reForInvalid := regexp.MustCompile(`\{\@for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*:=\s*range\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+trackBy\s+([a-zA-Z0-9_.]+)\}`)

	reEndFor := regexp.MustCompile(`\{\@endfor\}`)

	// Check for invalid syntax (missing index/underscore)
	if invalidMatches := reForInvalid.FindAllString(src, -1); len(invalidMatches) > 0 {
		lines := strings.Split(src, "\n")
		var invalidLines []int

		for i, line := range lines {
			if reForInvalid.MatchString(line) && !reFor.MatchString(line) {
				invalidLines = append(invalidLines, i+1)
			}
		}

		if len(invalidLines) > 0 {
			return "", fmt.Errorf("template syntax error in %s: Invalid {@for} syntax at line(s): %v\n"+
				"  The {@for} directive requires both index and value variables.\n"+
				"  Correct syntax: {@for index, value := range Slice trackBy value.Field}\n"+
				"  To ignore the index, use underscore: {@for _, value := range Slice trackBy value.Field}\n"+
				"  Example: {@for _, user := range Users trackBy user.ID}",
				templatePath, invalidLines)
		}
	}

	// Count directives to validate structure
	forCount := len(reFor.FindAllString(src, -1))
	endForCount := len(reEndFor.FindAllString(src, -1))

	if forCount != endForCount {
		// Find the line numbers to help the developer
		lines := strings.Split(src, "\n")
		var forLines []int
		var endForLines []int

		for i, line := range lines {
			if reFor.MatchString(line) {
				forLines = append(forLines, i+1)
			}
			if reEndFor.MatchString(line) {
				endForLines = append(endForLines, i+1)
			}
		}

		if forCount > endForCount {
			return "", fmt.Errorf("template validation error in %s: found %d {@for} directive(s) but only %d {@endfor} directive(s).\n"+
				"  {@for} found at line(s): %v\n"+
				"  {@endfor} found at line(s): %v\n"+
				"  Missing %d {@endfor} directive(s).",
				templatePath, forCount, endForCount, forLines, endForLines, forCount-endForCount)
		} else {
			return "", fmt.Errorf("template validation error in %s: found %d {@endfor} directive(s) but only %d {@for} directive(s).\n"+
				"  {@for} found at line(s): %v\n"+
				"  {@endfor} found at line(s): %v\n"+
				"  Extra %d {@endfor} directive(s) without matching {@for}.",
				templatePath, endForCount, forCount, forLines, endForLines, endForCount-forCount)
		}
	}

	// Transform {@for i, user := range Users trackBy user.ID} to placeholder elements
	src = reFor.ReplaceAllStringFunc(src, func(m string) string {
		matches := reFor.FindStringSubmatch(m)
		indexVar := matches[1]
		valueVar := matches[2]
		rangeExpr := matches[3]
		trackByExpr := matches[4]
		return fmt.Sprintf(`<go-for data-index="%s" data-value="%s" data-range="%s" data-trackby="%s">`,
			indexVar, valueVar, rangeExpr, trackByExpr)
	})

	src = reEndFor.ReplaceAllString(src, "</go-for>")
	return src, nil
}

// preprocessConditionals preprocesses template source to extract conditional blocks and replace them with placeholder nodes.
// It validates that every {@if} has a matching {@endif}.
func preprocessConditionals(src string, templatePath string) (string, error) {
	reIf := regexp.MustCompile(`\{\@if ([^}]+)\}`)
	reElseIf := regexp.MustCompile(`\{\@else if ([^}]+)\}`)
	reElse := regexp.MustCompile(`\{\@else\}`)
	reEndIf := regexp.MustCompile(`\{\@endif\}`)

	// Count directives to validate structure
	ifCount := len(reIf.FindAllString(src, -1))
	endifCount := len(reEndIf.FindAllString(src, -1))

	if ifCount != endifCount {
		// Find the line numbers to help the developer
		lines := strings.Split(src, "\n")
		var ifLines []int
		var endifLines []int

		for i, line := range lines {
			if reIf.MatchString(line) {
				ifLines = append(ifLines, i+1)
			}
			if reEndIf.MatchString(line) {
				endifLines = append(endifLines, i+1)
			}
		}

		if ifCount > endifCount {
			return "", fmt.Errorf("template validation error in %s: found %d {@if} directive(s) but only %d {@endif} directive(s).\n"+
				"  {@if} found at line(s): %v\n"+
				"  {@endif} found at line(s): %v\n"+
				"  Missing %d {@endif} directive(s).",
				templatePath, ifCount, endifCount, ifLines, endifLines, ifCount-endifCount)
		} else {
			return "", fmt.Errorf("template validation error in %s: found %d {@endif} directive(s) but only %d {@if} directive(s).\n"+
				"  {@if} found at line(s): %v\n"+
				"  {@endif} found at line(s): %v\n"+
				"  Extra %d {@endif} directive(s) without matching {@if}.",
				templatePath, endifCount, ifCount, ifLines, endifLines, endifCount-ifCount)
		}
	}

	src = reIf.ReplaceAllStringFunc(src, func(m string) string {
		cond := reIf.FindStringSubmatch(m)[1]
		return fmt.Sprintf("<go-conditional><go-if data-cond=\"%s\">", cond)
	})
	src = reElseIf.ReplaceAllStringFunc(src, func(m string) string {
		cond := reElseIf.FindStringSubmatch(m)[1]
		return fmt.Sprintf("</go-if><go-elseif data-cond=\"%s\">", cond)
	})
	// Handle {@else} - it closes the previous branch and opens go-else
	src = reElse.ReplaceAllString(src, func() string {
		// Check if the previous element is go-if or go-elseif
		// We need to close whichever was opened
		return "</go-if></go-elseif><go-else>"
	}())
	// {@endif} closes the last opened branch and the wrapper
	src = reEndIf.ReplaceAllString(src, "</go-if></go-elseif></go-else></go-conditional>")
	return src, nil
}

// estimateLineNumber tries to find the approximate line number where text appears in HTML source.
func estimateLineNumber(htmlSource, text string) int {
	lines := strings.Split(htmlSource, "\n")
	for i, line := range lines {
		if strings.Contains(line, text) {
			return i + 1 // Line numbers are 1-indexed
		}
	}
	return 1 // Default to line 1 if not found
}

// isBooleanAttribute checks if an attribute name is a standard HTML boolean attribute.
func isBooleanAttribute(attrName string) bool {
	return standardBooleanAttrs[attrName]
}

// validateBooleanCondition validates that a condition references a boolean field on the component.
// Returns the propertyDescriptor if valid, or exits with a compile error.
func validateBooleanCondition(condition string, comp componentInfo, templatePath string, lineNumber int, htmlSource string) propertyDescriptor {
	propDesc, exists := comp.Schema.Props[strings.ToLower(condition)]
	if !exists {
		contextLines := getContextLines(htmlSource, lineNumber, 2)
		availableFields := strings.Join(getAvailableFieldNames(comp.Schema.Props), ", ")
		fmt.Fprintf(os.Stderr, "Compilation Error in %s:%d: Condition '%s' not found on component '%s'. Available fields: [%s]\n%s",
			templatePath, lineNumber, condition, comp.PascalName, availableFields, contextLines)
		os.Exit(1)
	}
	if propDesc.GoType != "bool" {
		contextLines := getContextLines(htmlSource, lineNumber, 2)
		fmt.Fprintf(os.Stderr, "Compilation Error in %s:%d: Condition '%s' must be a bool field, found type '%s'.\n%s",
			templatePath, lineNumber, condition, propDesc.GoType, contextLines)
		os.Exit(1)
	}
	return propDesc
}

// generateTernaryExpression generates Go code for a ternary conditional expression.
// Supports negation operator: if negated is true, inverts the condition.
func generateTernaryExpression(negated bool, condition, trueVal, falseVal, receiver string, propDesc propertyDescriptor) string {
	if negated {
		// Swap true and false values for negation
		trueVal, falseVal = falseVal, trueVal
	}
	return fmt.Sprintf(`func() string {
		if %s.%s {
			return %s
		}
		return %s
	}()`, receiver, propDesc.Name, strconv.Quote(trueVal), strconv.Quote(falseVal))
}

// Compile is the main entry point for the AOT compiler.
func compile(srcDir, outDir string, devWarnings bool) error {
	opts := compileOptions{DevWarnings: devWarnings}

	// Step 1: Discover component templates and inspect their Go structs for props.
	components, err := discoverAndInspectComponents(srcDir)
	if err != nil {
		return fmt.Errorf("failed to discover or inspect components: %w", err)
	}
	fmt.Printf("Discovered and inspected %d component templates.\n", len(components))

	componentMap := make(map[string]componentInfo)
	for _, comp := range components {
		componentMap[comp.LowercaseName] = comp
	}

	// Step 2: Loop through each discovered component and compile its template.
	for _, comp := range components {
		if err := compileComponentTemplate(comp, componentMap, outDir, opts); err != nil {
			return fmt.Errorf("failed to compile template for %s: %w", comp.PascalName, err)
		}
	}
	return nil
}

// discoverAndInspectComponents finds all *.gt.html files and inspects their corresponding .go files.
func discoverAndInspectComponents(rootDir string) ([]componentInfo, error) {
	var components []componentInfo

	// Step 1: Load all packages in the module, configured for WASM.
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles, // Request file info
		Dir:  rootDir,
		Env:  append(os.Environ(), "GOOS=js", "GOARCH=wasm"),
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	// Step 2: Iterate through the loaded packages.
	for _, pkg := range pkgs {
		if len(pkg.GoFiles) == 0 {
			continue // Skip packages that are empty for the js/wasm target.
		}

		// All files in a package share the same directory.
		packageDir := filepath.Dir(pkg.GoFiles[0])

		// Step 3: Scan the package's directory for component templates (*.gt.html).
		files, err := os.ReadDir(packageDir)
		if err != nil {
			fmt.Printf("Warning: could not read directory %s: %v\n", packageDir, err)
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".gt.html") {
				continue
			}

			// We found a component template.
			templatePath := filepath.Join(packageDir, file.Name())
			pascalName := strings.TrimSuffix(file.Name(), ".gt.html")
			goFilePath := filepath.Join(packageDir, strings.ToLower(pascalName)+".go")

			schema, err := inspectGoFile(goFilePath, pascalName)
			if err != nil {
				fmt.Printf("Warning: could not inspect Go file %s: %v\n", goFilePath, err)
				schema = componentSchema{
					Props:   make(map[string]propertyDescriptor),
					Methods: make(map[string]bool),
				}
			}

			components = append(components, componentInfo{
				Path:          templatePath,
				PascalName:    pascalName,
				LowercaseName: strings.ToLower(pascalName),
				PackageName:   pkg.Name, // Use the package name from the loader.
				Schema:        schema,
			})
		}
	}

	if len(components) == 0 {
		fmt.Println("Warning: No component templates (*.gt.html) were found in any Go packages.")
	}

	return components, nil
}

// extractTypeName extracts the type name from an AST expression.
// Handles simple types (int, string, bool), slice types ([]User), and pointer types (*User).
func extractTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Simple type like "string", "int", "bool"
		return t.Name
	case *ast.ArrayType:
		// Slice or array type like "[]User"
		elemType := extractTypeName(t.Elt)
		return "[]" + elemType
	case *ast.StarExpr:
		// Pointer type like "*User"
		elemType := extractTypeName(t.X)
		return "*" + elemType
	case *ast.SelectorExpr:
		// Qualified type like "time.Time"
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name
		}
	}
	return "unknown"
}

// inspectStructInFile is a helper that inspects a specific struct type in a Go file.
// It returns a schema with the struct's exported fields.
func inspectStructInFile(path, structName string) (componentSchema, error) {
	schema := componentSchema{
		Props:   make(map[string]propertyDescriptor),
		Methods: make(map[string]bool),
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return schema, err
	}

	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				found = true
				for _, field := range structType.Fields.List {
					if len(field.Names) > 0 && field.Names[0].IsExported() {
						fieldName := field.Names[0].Name
						goType := extractTypeName(field.Type)
						schema.Props[strings.ToLower(fieldName)] = propertyDescriptor{
							Name:          fieldName,
							LowercaseName: strings.ToLower(fieldName),
							GoType:        goType,
						}
					}
				}
			}
		}
		return true
	})

	if !found {
		return schema, fmt.Errorf("struct '%s' not found in file", structName)
	}

	return schema, nil
}

// inspectGoFile parses a Go file and extracts the prop schema for a given struct.
func inspectGoFile(path, structName string) (componentSchema, error) {
	schema := componentSchema{
		Props:   make(map[string]propertyDescriptor),
		Methods: make(map[string]bool),
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return schema, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		// Inspect for struct fields (Props)
		if typeSpec, ok := n.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				for _, field := range structType.Fields.List {
					if len(field.Names) > 0 && field.Names[0].IsExported() {
						fieldName := field.Names[0].Name
						goType := extractTypeName(field.Type)
						schema.Props[strings.ToLower(fieldName)] = propertyDescriptor{
							Name:          fieldName,
							LowercaseName: strings.ToLower(fieldName),
							GoType:        goType,
						}
					}
				}
			}
		}

		// Inspect for methods (Event Handlers)
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			recv := funcDecl.Recv.List[0].Type
			if starExpr, ok := recv.(*ast.StarExpr); ok {
				recv = starExpr.X
			}
			if typeIdent, ok := recv.(*ast.Ident); ok && typeIdent.Name == structName {
				if funcDecl.Name.IsExported() {
					schema.Methods[funcDecl.Name.Name] = true
				}
			}
		}

		return true
	})

	return schema, nil
}

// compileComponentTemplate handles the code generation for a single component.
func compileComponentTemplate(comp componentInfo, componentMap map[string]componentInfo, outDir string, opts compileOptions) error {
	htmlContent, err := os.ReadFile(comp.Path)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", comp.Path, err)
	}
	htmlString := string(htmlContent)

	// Preprocess conditional blocks with validation
	htmlString, err = preprocessConditionals(htmlString, comp.Path)
	if err != nil {
		return err // Error message already includes template path and details
	}

	// Preprocess for-loop blocks with validation
	htmlString, err = preprocessFor(htmlString, comp.Path)
	if err != nil {
		return err // Error message already includes template path and details
	}

	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}
	bodyNode := findBody(doc)
	if bodyNode == nil {
		return fmt.Errorf("could not find <body> tag")
	}

	rootElement := findFirstElementChild(bodyNode)
	if rootElement == nil {
		return fmt.Errorf("no element found inside <body> tag to compile")
	}

	// Generate code for a single root node
	generatedCode := generateNodeCode(rootElement, "c", componentMap, comp, htmlString, opts, nil)

	template := `//go:build js || wasm
// +build js wasm

// Code generated by the nojs AOT compiler. DO NOT EDIT.
package %[2]s

import (
	"fmt"
	"github.com/vcrobe/nojs/vdom"
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/console"
	"strconv" // Added for type conversions
)

// Render generates the VNode tree for the %[1]s component.
func (c *%[1]s) Render(r *runtime.Renderer) *vdom.VNode {
	_ = strconv.Itoa // Suppress unused import error if no props are converted
	_ = fmt.Sprintf  // Suppress unused import error if no bindings are used
	_ = console.Log  // Suppress unused import error if no loops use dev warnings
	return %[3]s
}
`

	source := fmt.Sprintf(template, comp.PascalName, comp.PackageName, generatedCode)

	// Format the generated source code
	formattedSource, err := format.Source([]byte(source))
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	outFileName := fmt.Sprintf("%s.generated.go", comp.PascalName)
	outFilePath := filepath.Join(outDir, outFileName)
	return os.WriteFile(outFilePath, formattedSource, 0644)
}

// generateAttributesMap is a helper to create the Go map literal for an element's attributes.
func generateAttributesMap(n *html.Node, receiver string, currentComp componentInfo, htmlSource string) string {
	var attrs, events []string
	for _, a := range n.Attr {
		if after, ok := strings.CutPrefix(a.Key, "@"); ok {
			eventName := after
			handlerName := a.Val
			// Compile-time safety check!
			if _, ok := currentComp.Schema.Methods[handlerName]; !ok {
				// Find the line number for this event attribute
				lineNumber := findEventLineNumber(n, eventName, htmlSource)
				availableMethods := getAvailableMethodNames(currentComp.Schema.Methods)
				contextLines := getContextLines(htmlSource, lineNumber, 2)
				fmt.Fprintf(os.Stderr, "Compilation Error in %s:%d: Method '%s' not found on component '%s'. Available methods: [%s]\n%s",
					currentComp.Path, lineNumber, handlerName, currentComp.PascalName, availableMethods, contextLines)
				os.Exit(1)
			}
			switch eventName {
			case "onclick":
				// Generate the Go code to reference the component's method.
				handler := fmt.Sprintf(`%s.%s`, receiver, handlerName)
				events = append(events, fmt.Sprintf(`"onClick": %s`, handler))
			default:
				fmt.Printf("Warning: Unknown event directive '@%s' in %s.\n", eventName, currentComp.Path)
			}
		} else {
			// Check for inline conditional expressions in attribute values
			attrValue := a.Val
			lineNum := estimateLineNumber(htmlSource, fmt.Sprintf(`%s="%s"`, a.Key, attrValue))

			// Pattern 1: Boolean attribute shorthand {condition} or {!condition}
			if match := booleanShorthandRegex.FindStringSubmatch(attrValue); match != nil {
				negated := match[1] == "!"
				condition := match[2]

				// Only allow boolean shorthand for standard boolean attributes
				if !isBooleanAttribute(a.Key) {
					contextLines := getContextLines(htmlSource, lineNum, 2)
					fmt.Fprintf(os.Stderr, "Compilation Error in %s:%d: Boolean shorthand syntax can only be used with standard HTML boolean attributes. For attribute '%s', use the full ternary expression: {%s ? 'true' : 'false'}\n%s",
						currentComp.Path, lineNum, a.Key, condition, contextLines)
					os.Exit(1)
				}

				// Validate condition is a boolean field
				propDesc := validateBooleanCondition(condition, currentComp, currentComp.Path, lineNum, htmlSource)

				// Generate conditional code: if negated, invert the condition
				if negated {
					attrs = append(attrs, fmt.Sprintf(`"%s": !%s.%s`, a.Key, receiver, propDesc.Name))
				} else {
					attrs = append(attrs, fmt.Sprintf(`"%s": %s.%s`, a.Key, receiver, propDesc.Name))
				}
				continue
			}

			// Pattern 2: Ternary expressions in attribute values
			if ternaryExprRegex.MatchString(attrValue) {
				// Replace all ternary expressions in the value
				result := attrValue
				ternaryMatches := ternaryExprRegex.FindAllStringSubmatch(attrValue, -1)

				for _, match := range ternaryMatches {
					fullMatch := match[0]
					negated := match[1] == "!"
					condition := match[2]
					trueVal := match[3]
					falseVal := match[4]

					// Validate condition is a boolean field
					propDesc := validateBooleanCondition(condition, currentComp, currentComp.Path, lineNum, htmlSource)

					// Generate ternary expression
					ternaryCode := generateTernaryExpression(negated, condition, trueVal, falseVal, receiver, propDesc)

					// If the attribute value is only the ternary expression
					if result == fullMatch {
						attrs = append(attrs, fmt.Sprintf(`"%s": %s`, a.Key, ternaryCode))
						result = ""
						break
					}

					// Otherwise, replace in the string (for concatenation)
					result = strings.Replace(result, fullMatch, "%s", 1)
				}

				// If there were other parts, wrap in fmt.Sprintf
				if result != "" {
					var args []string
					for _, match := range ternaryMatches {
						negated := match[1] == "!"
						condition := match[2]
						trueVal := match[3]
						falseVal := match[4]
						propDesc := validateBooleanCondition(condition, currentComp, currentComp.Path, lineNum, htmlSource)
						args = append(args, generateTernaryExpression(negated, condition, trueVal, falseVal, receiver, propDesc))
					}
					attrs = append(attrs, fmt.Sprintf(`"%s": fmt.Sprintf(%s, %s)`, a.Key, strconv.Quote(result), strings.Join(args, ", ")))
				}
				continue
			}

			// Pattern 3: Regular static attribute
			attrs = append(attrs, fmt.Sprintf(`"%s": "%s"`, a.Key, a.Val))
		}
	}

	if len(attrs) == 0 && len(events) == 0 {
		return "nil"
	}
	allProps := append(attrs, events...)
	return fmt.Sprintf("map[string]any{%s}", strings.Join(allProps, ", "))
}

// generateTextExpression handles data binding in text nodes.
// loopCtx can be nil if not inside a loop.
func generateTextExpression(text string, receiver string, currentComp componentInfo, htmlSource string, lineNumber int, loopCtx *loopContext) string {
	// Check for ternary expressions first
	ternaryMatches := ternaryExprRegex.FindAllStringSubmatch(text, -1)

	if len(ternaryMatches) > 0 {
		// Handle ternary expressions
		result := text

		for _, match := range ternaryMatches {
			fullMatch := match[0]
			negated := match[1] == "!"
			condition := match[2]
			trueVal := match[3]
			falseVal := match[4]

			// Validate condition is a boolean field
			propDesc := validateBooleanCondition(condition, currentComp, currentComp.Path, lineNumber, htmlSource)

			// Generate ternary expression
			ternaryCode := generateTernaryExpression(negated, condition, trueVal, falseVal, receiver, propDesc)

			// If the text contains only the ternary expression, return it directly
			if result == fullMatch {
				return ternaryCode
			}

			// Otherwise, replace the match with a placeholder for fmt.Sprintf
			result = strings.Replace(result, fullMatch, "%s", 1)
		}

		// If there are other parts of the text, wrap in fmt.Sprintf
		var args []string
		for _, match := range ternaryMatches {
			negated := match[1] == "!"
			condition := match[2]
			trueVal := match[3]
			falseVal := match[4]
			propDesc := validateBooleanCondition(condition, currentComp, currentComp.Path, lineNumber, htmlSource)
			args = append(args, generateTernaryExpression(negated, condition, trueVal, falseVal, receiver, propDesc))
		}

		return fmt.Sprintf(`fmt.Sprintf(%s, %s)`, strconv.Quote(result), strings.Join(args, ", "))
	}

	// Original data binding logic
	matches := dataBindingRegex.FindAllStringSubmatch(text, -1)

	if len(matches) == 0 {
		return strconv.Quote(text) // It's just a static string
	}

	formatString := dataBindingRegex.ReplaceAllString(text, "%v")
	var args []string

	for _, match := range matches {
		fieldName := match[1]

		// Check if this is a loop variable first
		if loopCtx != nil {
			if fieldName == loopCtx.IndexVar {
				// Reference loop index variable
				args = append(args, fieldName)
				continue
			}
			if fieldName == loopCtx.ValueVar {
				// Reference loop value variable
				args = append(args, fieldName)
				continue
			}
			// Check if it's a field access on the loop value variable (e.g., user.Name)
			if strings.Contains(fieldName, ".") {
				parts := strings.SplitN(fieldName, ".", 2)
				varName := parts[0]
				if varName == loopCtx.ValueVar {
					// This is a field access on the loop value variable
					// Just use it as-is (e.g., user.Name)
					args = append(args, fieldName)
					continue
				}
			}
		}

		// Type-safety check: does the field exist on the component struct?
		if _, ok := currentComp.Schema.Props[strings.ToLower(fieldName)]; !ok {
			// If we're in a loop, provide more context in the error
			if loopCtx != nil {
				fmt.Fprintf(os.Stderr, "Compilation Error in %s: Field '%s' not found.\n"+
					"  - Not a loop variable (loop has: %s, %s)\n"+
					"  - Not a component field (available: %s)\n"+
					"  - For loop item fields, use: %s.FieldName\n",
					currentComp.Path, fieldName,
					loopCtx.IndexVar, loopCtx.ValueVar,
					strings.Join(getAvailableFieldNames(currentComp.Schema.Props), ", "),
					loopCtx.ValueVar)
			} else {
				fmt.Fprintf(os.Stderr, "Compilation Error in %s: Field '%s' not found on component '%s' for data binding.\n",
					currentComp.Path, fieldName, currentComp.PascalName)
			}
			os.Exit(1)
		}
		args = append(args, fmt.Sprintf("%s.%s", receiver, fieldName))
	}

	return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, formatString, strings.Join(args, ", "))
}

// generateForLoopCode generates Go for...range loop code for list rendering.
func generateForLoopCode(n *html.Node, receiver string, componentMap map[string]componentInfo, currentComp componentInfo, htmlSource string, opts compileOptions) string {
	// Extract loop variables from data attributes
	indexVar := ""
	valueVar := ""
	rangeExpr := ""
	trackByExpr := ""

	for _, attr := range n.Attr {
		switch attr.Key {
		case "data-index":
			indexVar = attr.Val
		case "data-value":
			valueVar = attr.Val
		case "data-range":
			rangeExpr = attr.Val
		case "data-trackby":
			trackByExpr = attr.Val
		}
	}

	// Validate that we have the required attributes
	if valueVar == "" || rangeExpr == "" || trackByExpr == "" {
		fmt.Fprintf(os.Stderr, "Compilation Error in %s: Invalid {@for} directive - missing required attributes.\n", currentComp.Path)
		os.Exit(1)
	}

	// Validate that the range expression exists on the component
	propDesc, exists := currentComp.Schema.Props[strings.ToLower(rangeExpr)]
	if !exists {
		availableFields := strings.Join(getAvailableFieldNames(currentComp.Schema.Props), ", ")
		fmt.Fprintf(os.Stderr, "Compilation Error in %s: Field '%s' not found on component '%s'. Available fields: [%s]\n",
			currentComp.Path, rangeExpr, currentComp.PascalName, availableFields)
		os.Exit(1)
	}

	// Validate that the field is a slice type
	if !strings.HasPrefix(propDesc.GoType, "[]") {
		fmt.Fprintf(os.Stderr, "Compilation Error in %s: Field '%s' must be a slice or array type for {@for} directive, found type '%s'.\n",
			currentComp.Path, rangeExpr, propDesc.GoType)
		os.Exit(1)
	}

	// Validate trackBy expression
	// Parse trackBy to extract variable and field: "user.ID" -> variable="user", field="ID"
	trackByParts := strings.Split(trackByExpr, ".")
	if len(trackByParts) != 2 {
		fmt.Fprintf(os.Stderr, "Compilation Error in %s: trackBy expression '%s' must be in format 'variable.Field' (e.g., 'user.ID').\n",
			currentComp.Path, trackByExpr)
		os.Exit(1)
	}

	trackByVar := trackByParts[0]
	trackByField := trackByParts[1]

	// Verify the variable matches the loop value variable
	if trackByVar != valueVar {
		fmt.Fprintf(os.Stderr, "Compilation Error in %s: trackBy variable '%s' must match the loop value variable '%s'.\n"+
			"  Expected: trackBy %s.FieldName\n",
			currentComp.Path, trackByVar, valueVar, valueVar)
		os.Exit(1)
	}

	// Extract element type from slice type: "[]User" -> "User"
	elementType := strings.TrimPrefix(propDesc.GoType, "[]")

	// Validate that the trackBy field exists on the element type
	// We need to inspect the element type's struct definition
	goFilePath := filepath.Join(filepath.Dir(currentComp.Path), strings.ToLower(currentComp.PascalName)+".go")
	elementSchema, err := inspectStructInFile(goFilePath, elementType)
	if err != nil {
		// If we can't find the struct in the component file, it might be defined elsewhere
		// For now, we'll skip validation with a warning
		fmt.Fprintf(os.Stderr, "Warning in %s: Could not validate trackBy field '%s' on type '%s': %v\n",
			currentComp.Path, trackByField, elementType, err)
	} else {
		// Check if the trackBy field exists on the element type (case-insensitive lookup)
		propDesc, exists := elementSchema.Props[strings.ToLower(trackByField)]
		if !exists {
			availableFields := strings.Join(getAvailableFieldNames(elementSchema.Props), ", ")
			fmt.Fprintf(os.Stderr, "Compilation Error in %s: trackBy identifier '%s' not found on type '%s'.\nAvailable fields: [%s]\n",
				currentComp.Path, trackByField, elementType, availableFields)
			os.Exit(1)
		}

		// Verify exact case match - the field name in the template must match the actual struct field
		if propDesc.Name != trackByField {
			availableFields := strings.Join(getAvailableFieldNames(elementSchema.Props), ", ")
			fmt.Fprintf(os.Stderr, "Compilation Error in %s: trackBy identifier '%s' not found on type '%s'.\nAvailable fields: [%s]\n",
				currentComp.Path, trackByField, elementType, availableFields)
			os.Exit(1)
		}
	}

	// Generate the loop body - collect child VNodes
	var code strings.Builder

	// Generate IIFE that returns a slice of VNodes
	code.WriteString("func() []*vdom.VNode {\n")
	code.WriteString(fmt.Sprintf("\tvar %s_nodes []*vdom.VNode\n", valueVar))

	// Add development warning if enabled
	if opts.DevWarnings {
		code.WriteString(fmt.Sprintf("\t// Development warning for empty slice\n"))
		code.WriteString(fmt.Sprintf("\tif len(%s.%s) == 0 {\n", receiver, propDesc.Name))
		code.WriteString(fmt.Sprintf("\t\tconsole.Warning(\"[@for] Rendering empty list for '%s' in %s. Consider using {@if} to handle empty state.\")\n",
			propDesc.Name, currentComp.PascalName))
		code.WriteString("\t}\n\n")
	}

	// Generate the for loop
	code.WriteString(fmt.Sprintf("\tfor %s, %s := range %s.%s {\n", indexVar, valueVar, receiver, propDesc.Name))

	// Create loop context for child nodes
	loopCtx := &loopContext{
		IndexVar: indexVar,
		ValueVar: valueVar,
	}

	// Generate code for each child node in the loop body
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode || (c.Type == html.TextNode && strings.TrimSpace(c.Data) != "") {
			childCode := generateNodeCode(c, receiver, componentMap, currentComp, htmlSource, opts, loopCtx)
			if childCode != "" {
				code.WriteString(fmt.Sprintf("\t\t%s_child := %s\n", valueVar, childCode))
				code.WriteString(fmt.Sprintf("\t\tif %s_child != nil {\n", valueVar))
				code.WriteString(fmt.Sprintf("\t\t\t%s_nodes = append(%s_nodes, %s_child)\n", valueVar, valueVar, valueVar))
				code.WriteString("\t\t}\n")
			}
		}
	}

	code.WriteString("\t}\n")
	code.WriteString(fmt.Sprintf("\treturn %s_nodes\n", valueVar))
	code.WriteString("}()")

	return code.String()
}

// generateConditionalCode generates Go if/else blocks for conditional rendering.
func generateConditionalCode(n *html.Node, receiver string, componentMap map[string]componentInfo, currentComp componentInfo, htmlSource string, opts compileOptions, loopCtx *loopContext) string {
	var code strings.Builder

	// Generate IIFE (Immediately Invoked Function Expression)
	code.WriteString("func() *vdom.VNode {\n")

	// Process children of go-conditional wrapper
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "go-if" {
			// Extract and validate condition
			cond := ""
			for _, attr := range c.Attr {
				if attr.Key == "data-cond" {
					cond = attr.Val
					break
				}
			}

			propDesc, exists := currentComp.Schema.Props[strings.ToLower(cond)]
			if !exists {
				fmt.Fprintf(os.Stderr, "Compilation Error in %s: Condition '%s' not found on component '%s'.\n", currentComp.Path, cond, currentComp.PascalName)
				os.Exit(1)
			}
			if propDesc.GoType != "bool" {
				fmt.Fprintf(os.Stderr, "Compilation Error in %s: Condition '%s' must be a bool field, found type '%s'.\n", currentComp.Path, cond, propDesc.GoType)
				os.Exit(1)
			}

			code.WriteString(fmt.Sprintf("if %s.%s {\n", receiver, propDesc.Name))
			foundContent := false
			for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
				childCode := generateNodeCode(cc, receiver, componentMap, currentComp, htmlSource, opts, loopCtx)
				if childCode != "" {
					code.WriteString("return ")
					code.WriteString(childCode)
					code.WriteString("\n")
					foundContent = true
					break
				}
			}
			if !foundContent {
				code.WriteString("return nil\n")
			}
			code.WriteString("}")
		} else if c.Type == html.ElementNode && c.Data == "go-elseif" {
			// Extract and validate condition
			elseifCond := ""
			for _, attr := range c.Attr {
				if attr.Key == "data-cond" {
					elseifCond = attr.Val
					break
				}
			}

			propDesc, exists := currentComp.Schema.Props[strings.ToLower(elseifCond)]
			if !exists {
				fmt.Fprintf(os.Stderr, "Compilation Error in %s: Condition '%s' not found on component '%s'.\n", currentComp.Path, elseifCond, currentComp.PascalName)
				os.Exit(1)
			}
			if propDesc.GoType != "bool" {
				fmt.Fprintf(os.Stderr, "Compilation Error in %s: Condition '%s' must be a bool field, found type '%s'.\n", currentComp.Path, elseifCond, propDesc.GoType)
				os.Exit(1)
			}

			code.WriteString(fmt.Sprintf(" else if %s.%s {\n", receiver, propDesc.Name))
			foundContent := false
			for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
				childCode := generateNodeCode(cc, receiver, componentMap, currentComp, htmlSource, opts, loopCtx)
				if childCode != "" {
					code.WriteString("return ")
					code.WriteString(childCode)
					code.WriteString("\n")
					foundContent = true
					break
				}
			}
			if !foundContent {
				code.WriteString("return nil\n")
			}
			code.WriteString("}")
		} else if c.Type == html.ElementNode && c.Data == "go-else" {
			code.WriteString(" else {\n")
			foundContent := false
			for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
				childCode := generateNodeCode(cc, receiver, componentMap, currentComp, htmlSource, opts, loopCtx)
				if childCode != "" {
					code.WriteString("return ")
					code.WriteString(childCode)
					code.WriteString("\n")
					foundContent = true
					break
				}
			}
			if !foundContent {
				code.WriteString("return nil\n")
			}
			code.WriteString("}\n")
			// Don't add the fallback return nil after else block
			code.WriteString("}()")
			return code.String()
		}
	}

	// Only add fallback return nil if there's no else branch
	code.WriteString("\nreturn nil\n}()")
	return code.String()
}

// generateNodeCode recursively generates Go vdom calls.
// loopCtx can be nil if not inside a loop.
func generateNodeCode(n *html.Node, receiver string, componentMap map[string]componentInfo, currentComp componentInfo, htmlSource string, opts compileOptions, loopCtx *loopContext) string {
	if n.Type == html.TextNode {
		content := strings.TrimSpace(n.Data)
		if content == "" {
			return ""
		}
		// In a real scenario, you'd handle text content more robustly.
		// For now, we assume text is primarily for simple elements like <p>.
		return ""
	}

	if n.Type == html.ElementNode {
		tagName := n.Data

		// 0. Handle conditional placeholder nodes
		if tagName == "go-conditional" {
			return generateConditionalCode(n, receiver, componentMap, currentComp, htmlSource, opts, loopCtx)
		}
		if tagName == "go-if" || tagName == "go-elseif" || tagName == "go-else" {
			// These are handled within go-conditional processing
			return ""
		}

		// 0.5. Handle for-loop placeholder nodes
		if tagName == "go-for" {
			return generateForLoopCode(n, receiver, componentMap, currentComp, htmlSource, opts)
		}

		// 1. Handle Custom Components
		if compInfo, isComponent := componentMap[tagName]; isComponent {
			propsStr := generateStructLiteral(n, compInfo, htmlSource, currentComp.Path)
			key := fmt.Sprintf("%s_%d", compInfo.PascalName, childCount(n.Parent, n)) // Simple key generation

			return fmt.Sprintf(`r.RenderChild("%s", &%s%s)`, key, compInfo.PascalName, propsStr)
		}

		// 2. Handle Standard HTML Elements
		var childrenCode []string
		hasForLoop := false
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			// Check if this child is a go-for node
			if c.Type == html.ElementNode && c.Data == "go-for" {
				hasForLoop = true
			}
			childCode := generateNodeCode(c, receiver, componentMap, currentComp, htmlSource, opts, loopCtx)
			if childCode != "" {
				childrenCode = append(childrenCode, childCode)
			}
		}

		var childrenStr string
		if hasForLoop {
			// When we have a for loop, we need to build children differently
			// Generate code that collects all children into a slice
			childrenStr = "func() []*vdom.VNode {\nvar allChildren []*vdom.VNode\n"
			for _, code := range childrenCode {
				// Check if this looks like a for loop return (starts with "func")
				if strings.HasPrefix(strings.TrimSpace(code), "func()") {
					childrenStr += fmt.Sprintf("allChildren = append(allChildren, %s...)\n", code)
				} else {
					childrenStr += fmt.Sprintf("allChildren = append(allChildren, %s)\n", code)
				}
			}
			childrenStr += "return allChildren\n}()..."
		} else {
			childrenStr = strings.Join(childrenCode, ", ")
		}

		attrsMapStr := generateAttributesMap(n, receiver, currentComp, htmlSource)

		switch tagName {
		case "div", "ul", "ol":
			return fmt.Sprintf("vdom.Div(%s, %s)", attrsMapStr, childrenStr)
		case "p", "button", "li", "h1", "h2", "h3", "h4", "h5", "h6":
			textContent := ""
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				// Handle data binding and inline conditionals in the text content
				// Estimate line number by searching for the text in the HTML source
				lineNum := estimateLineNumber(htmlSource, n.FirstChild.Data)
				textContent = generateTextExpression(n.FirstChild.Data, receiver, currentComp, htmlSource, lineNum, loopCtx)
			} else {
				textContent = `""` // Default to empty string if no text node
			}

			// The VDOM helpers expect a string, so we pass the generated expression
			switch tagName {
			case "p":
				return fmt.Sprintf("vdom.Paragraph(%s, %s)", textContent, attrsMapStr)
			case "button":
				return fmt.Sprintf("vdom.Button(%s, %s, %s)", textContent, attrsMapStr, childrenStr)
			default:
				// For li, h1-h6, use NewVNode directly with text content
				return fmt.Sprintf("vdom.NewVNode(%s, %s, nil, %s)", strconv.Quote(tagName), attrsMapStr, textContent)
			}
		default:
			return `vdom.Div(nil)` // Default to an empty div for unknown tags
		}
	}

	return ""
}

// generateStructLiteral creates the { Field: value, ... } string.
func generateStructLiteral(n *html.Node, compInfo componentInfo, htmlSource string, templatePath string) string {
	var props []string

	// Extract the original attribute names from the HTML source
	originalAttrs, lineNumber := extractOriginalAttributesWithLineNumber(n, compInfo.LowercaseName, htmlSource)

	for _, attr := range n.Attr {
		// Get the original casing from the source
		originalKey := attr.Key
		if origName, found := originalAttrs[attr.Key]; found {
			originalKey = origName
		}

		// Check if the ORIGINAL attribute starts with a capital letter
		if len(originalKey) > 0 && originalKey[0] >= 'A' && originalKey[0] <= 'Z' {
			// This is a prop binding - it must match an exported field
			lookupKey := strings.ToLower(originalKey)

			if propDesc, ok := compInfo.Schema.Props[lookupKey]; ok {
				valueStr := convertPropValue(attr.Val, propDesc.GoType)
				props = append(props, fmt.Sprintf("%s: %s", propDesc.Name, valueStr))
			} else {
				// Attribute starts with capital letter but doesn't match any exported field
				availableFields := strings.Join(getAvailableFieldNames(compInfo.Schema.Props), ", ")
				contextLines := getContextLines(htmlSource, lineNumber, 2)
				fmt.Fprintf(os.Stderr, "Compilation Error in %s:%d: Attribute '%s' does not match any exported field on component '%s'. Available fields: [%s]\n%s",
					templatePath, lineNumber, originalKey, compInfo.PascalName, availableFields, contextLines)
				os.Exit(1)
			}
		} else if propDesc, ok := compInfo.Schema.Props[attr.Key]; ok {
			// Lowercase attribute that happens to match a field
			valueStr := convertPropValue(attr.Val, propDesc.GoType)
			props = append(props, fmt.Sprintf("%s: %s", propDesc.Name, valueStr))
		}
	}

	if len(props) == 0 {
		return "{}"
	}

	return fmt.Sprintf("{%s}", strings.Join(props, ", "))
}

// extractOriginalAttributesWithLineNumber extracts the original attribute names and line number from the HTML source.
// This is needed because the HTML parser lowercases all attributes.
func extractOriginalAttributesWithLineNumber(n *html.Node, componentName, htmlSource string) (map[string]string, int) {
	originalAttrs := make(map[string]string)
	lineNumber := 1

	// Find the component tag in the HTML source (case-insensitive tag name)
	// Pattern: <componentName attr1="..." attr2="..." ...>
	pattern := fmt.Sprintf(`(?i)<%s\s+([^>]*)>`, regexp.QuoteMeta(componentName))
	re := regexp.MustCompile(pattern)
	matchIndex := re.FindStringSubmatchIndex(htmlSource)

	if matchIndex == nil || len(matchIndex) < 4 {
		return originalAttrs, lineNumber
	}

	// Calculate line number by counting newlines before the match
	lineNumber = strings.Count(htmlSource[:matchIndex[0]], "\n") + 1

	// Extract the attribute string
	attrString := htmlSource[matchIndex[2]:matchIndex[3]]

	// Extract individual attributes with their original casing
	// Pattern: attrName="value" or attrName='value'
	attrPattern := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s*=\s*["']([^"']*)["']`)
	attrMatches := attrPattern.FindAllStringSubmatch(attrString, -1)

	for _, match := range attrMatches {
		if len(match) >= 2 {
			originalName := match[1]
			lowercaseName := strings.ToLower(originalName)
			originalAttrs[lowercaseName] = originalName
		}
	}

	return originalAttrs, lineNumber
}

// getAvailableFieldNames returns a slice of exported field names for error messages.
func getAvailableFieldNames(props map[string]propertyDescriptor) []string {
	var names []string
	for _, prop := range props {
		names = append(names, prop.Name)
	}
	return names
}

// getAvailableMethodNames returns a comma-separated string of available method names for error messages.
func getAvailableMethodNames(methods map[string]bool) string {
	var names []string
	for methodName := range methods {
		names = append(names, methodName)
	}
	return strings.Join(names, ", ")
}

// findEventLineNumber finds the line number where an event attribute is defined.
func findEventLineNumber(n *html.Node, eventName, htmlSource string) int {
	// Look for the event attribute pattern: @eventName="..."
	// We need to find the element's tag and then the specific event attribute
	tagName := n.Data

	// Create a pattern to find this specific element with the event attribute
	// This is a simplified approach - it finds the first occurrence
	pattern := fmt.Sprintf(`(?i)<%s[^>]*@%s\s*=`, regexp.QuoteMeta(tagName), regexp.QuoteMeta(eventName))
	re := regexp.MustCompile(pattern)
	matchIndex := re.FindStringIndex(htmlSource)

	if matchIndex == nil {
		return 1 // Default to line 1 if not found
	}

	// Count newlines before the match to get the line number
	lineNumber := strings.Count(htmlSource[:matchIndex[0]], "\n") + 1
	return lineNumber
}

// getContextLines returns a formatted string with context lines around the error line.
// It shows 'contextSize' lines before and after the target line.
func getContextLines(source string, lineNumber int, contextSize int) string {
	lines := strings.Split(source, "\n")

	// Calculate the range of lines to show
	startLine := lineNumber - contextSize - 1 // -1 for 0-based indexing
	if startLine < 0 {
		startLine = 0
	}

	endLine := lineNumber + contextSize // lineNumber is already the index we want to highlight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	var result strings.Builder
	result.WriteString("\n")

	for i := startLine; i < endLine; i++ {
		lineNum := i + 1
		prefix := "  "

		// Highlight the error line with a marker
		if lineNum == lineNumber {
			prefix = "> "
		}

		result.WriteString(fmt.Sprintf("%s%4d | %s\n", prefix, lineNum, lines[i]))
	}

	return result.String()
}

// convertPropValue generates the Go code to convert a string to the target type.
func convertPropValue(value, goType string) string {
	switch goType {
	case "string":
		return strconv.Quote(value)
	case "int":
		// In a real compiler, you'd handle the error. Here we assume valid input.
		return fmt.Sprintf("func() int { i, _ := strconv.Atoi(\"%s\"); return i }()", value)
	case "bool":
		return fmt.Sprintf("func() bool { b, _ := strconv.ParseBool(\"%s\"); return b }()", value)
	default:
		// Default to string for unknown types
		return strconv.Quote(value)
	}
}

// findBody finds the <body> node in the parsed HTML.
func findBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findBody(c); result != nil {
			return result
		}
	}

	return nil
}

// findFirstElementChild finds the first actual element inside a node.
func findFirstElementChild(n *html.Node) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return c
		}
	}
	return nil
}

// childCount is a helper function to count preceding element siblings for key generation.
func childCount(parent *html.Node, until *html.Node) int {
	count := 0

	if parent == nil {
		return 0
	}

	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c == until {
			break
		}

		if c.Type == html.ElementNode {
			count++
		}
	}

	return count
}
