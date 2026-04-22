package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/norowachi/env-doctor/internal/parser"
	"github.com/norowachi/env-doctor/internal/reporter"
	"github.com/norowachi/env-doctor/internal/validator"
)

const version = "0.2.1"

func main() {
	_, JSON := os.LookupEnv("ENV_DOCTOR_JSON")

	examplePath := flag.String("example", ".env.example", "Path to the example/schema env file")
	envPath := flag.String("env", ".env", "Path to the target env file to check")
	jsonOutput := flag.Bool("json", JSON, "Output results as JSON (for CI pipelines)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	ignoreWarnings := flag.Bool("ignore-warnings", false, "Suppress warnings and info messages from output")
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
	if *ignoreWarnings {
		issues = filterWarnings(issues)
	}

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

func filterWarnings(issues []validator.Issue) []validator.Issue {
	filtered := issues[:0]
	for _, i := range issues {
		if i.Severity == validator.SeverityError {
			filtered = append(filtered, i)
		}
	}
	return filtered
}
