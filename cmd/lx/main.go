package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andre-carbajal/lx/internal/detector"
	"github.com/andre-carbajal/lx/internal/executor"
	"github.com/andre-carbajal/lx/internal/parser"
	"github.com/andre-carbajal/lx/internal/translator"
)

var Version = "dev"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		showVersion()
		return
	}

	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		showHelp()
		return
	}

	globalFlags, cmdArgs := parser.ExtractGlobalFlags(args)

	if len(cmdArgs) == 0 {
		showVersion()
		return
	}

	linuxCmd, err := parser.Parse(strings.Join(cmdArgs, " "))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse command: %v\n", err)
		os.Exit(1)
	}

	processProvider := detector.DefaultProvider()
	shellDetector := detector.New(processProvider)
	shell := shellDetector.Detect()

	router, err := translator.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize translator: %v\n", err)
		os.Exit(1)
	}

	translationResult, err := router.Translate(linuxCmd, shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Translation failed: %v\n", err)
		os.Exit(1)
	}

	for _, warning := range translationResult.Warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
	}

	exec := &executor.Executor{}
	exitCode := exec.Execute(translationResult.Translated, shell, globalFlags)
	os.Exit(exitCode)
}

func showVersion() {
	fmt.Printf("lx v%s — Linux-to-Windows Command Translator\n", Version)
	fmt.Println("Use 'lx --help' for usage information")
}

func showHelp() {
	fmt.Printf("lx v%s — Linux-to-Windows Command Translator\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  lx [OPTIONS] <COMMAND>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --dry        Show what command would be executed, but don't run it")
	fmt.Println("  --verbose    Show shell type and translated command before executing")
	fmt.Println("  -v           Verbose mode (alias for --verbose)")
	fmt.Println("  --help       Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lx ls -la")
	fmt.Println("  lx --dry cp file1.txt file2.txt")
	fmt.Println("  lx --verbose grep 'pattern' file.txt")
	fmt.Println()
	fmt.Println("Supported commands:")
	fmt.Println("  File operations: ls, cd, pwd, pushd, popd, cp, mv, rm, mkdir, touch, find")
	fmt.Println("  Text processing: cat, grep, head, tail, sort, echo, wc")
	fmt.Println("  System: clear, date, env, ps, kill, ping, curl, ipconfig")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  LX_SHELL     Override detected shell (cmd or powershell)")
}
