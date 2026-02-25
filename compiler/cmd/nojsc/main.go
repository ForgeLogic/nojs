package main

import (
	"flag"
	"fmt"
	"log"

	compiler "github.com/ForgeLogic/nojs-compiler"
)

func main() {
	inDir := flag.String("in", ".", "The source directory to scan for *.gt.html files.")
	devMode := flag.Bool("dev", false, "Enable development mode (warnings, verbose errors, panic on lifecycle failures)")
	flag.Parse()

	fmt.Printf("Starting compilation...\nSource directory: %s\n", *inDir)
	if *devMode {
		fmt.Printf("Development mode: ENABLED\n")
	}
	err := compiler.Compile(*inDir, *devMode)
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Printf("🎉 Compilation completed successfully!\n")
}
