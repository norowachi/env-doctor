package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/norowachi/env-doctor/internal/parser"
	"github.com/norowachi/env-doctor/internal/reporter"
	"github.com/norowachi/env-doctor/internal/validator"
)

const version = "0.1.0"

func main() {
	examplePath := flag.String("example", ".env.example", "Path to the example/schema env file")
	envPath := flag.String("env", ".env", "Path to the target env file to check")
	jsonOutput := flag.Bool("json", false, "Output results as JSON (for CI pipelines)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("env-doctor v%s\n", version)
		os.Exit(0)
	}

	// Parse env example/schema file
	example, err := parser.Parse(*examplePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", *examplePath, err)
		os.Exit(2)
	}

	// Parse target env file
	actual, err := parser.Parse(*envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", *envPath, err)
		os.Exit(2)
	}

	// Validate
	issues := validator.Validate(example, actual)

	// Generate Report
	if *jsonOutput {
		reporter.PrintJSON(issues)
		if hasErrors(issues) {
			os.Exit(1)
		}
		return
	}

	hasErrors := reporter.Print(issues, *examplePath, *envPath)
	if hasErrors {
		os.Exit(1)
	}
}

func hasErrors(issues []validator.Issue) bool {
	for _, i := range issues {
		if i.Severity == validator.SeverityError {
			return true
		}
	}
	return false
}
