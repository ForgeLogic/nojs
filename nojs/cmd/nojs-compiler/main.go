package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vcrobe/nojs/compiler"
)

func main() {
	// --- CLI Flags Updated for Directory-Based Compilation ---
	// The '-in' flag now specifies the source directory to scan for components.
	inDir := flag.String("in", ".", "The source directory to scan for *.gt.html files.")
	// The '-out' flag now specifies the directory where generated Go files will be placed.
	outDir := flag.String("out", "", "The output directory for the generated Go files.")
	// The '-dev' flag enables development mode (warnings, verbose errors, panic on lifecycle failures).
	devMode := flag.Bool("dev", false, "Enable development mode (warnings, verbose errors, panic on lifecycle failures)")
	flag.Parse()

	if *outDir == "" {
		log.Fatal("Error: The -out flag is required to specify the output directory.")
	}

	// Create the output directory if it doesn't exist.
	if err := os.MkdirAll(*outDir, 0755); err != nil {
		log.Fatalf("Error: Could not create output directory %s: %v", *outDir, err)
	}

	// The CLI's job is now to pass the directories to the core compiler logic.
	fmt.Printf("Starting compilation...\nSource directory: %s\nOutput directory: %s\n", *inDir, *outDir)
	if *devMode {
		fmt.Printf("Development mode: ENABLED\n")
	}
	err := compiler.Compile(*inDir, *outDir, *devMode)
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	// Success! Let the user know.
	fmt.Printf("🎉 Compilation completed successfully!\n")
}
