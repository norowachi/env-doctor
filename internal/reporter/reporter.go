package reporter

import (
	"fmt"
	"strings"

	"github.com/norowachi/env-doctor/internal/validator"
)

// color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// Render to stdout with color and final summary
func Print(issues []validator.Issue, examplePath, envPath string) (hasErrors bool) {
	fmt.Printf("\n%s🩺 env-doctor%s  %s→%s  %s vs %s\n",
		colorBold, colorReset, colorGray, colorReset, examplePath, envPath)
	fmt.Println(strings.Repeat("─", 55))

	if len(issues) == 0 {
		fmt.Printf("%s✓ All good! Your env matches the schema perfectly.%s\n\n", colorGreen, colorReset)
		return false
	}

	errors, warnings, infos := 0, 0, 0

	for _, issue := range issues {
		switch issue.Severity {
		case validator.SeverityError:
			fmt.Printf("%s✗ ERROR%s   %s%-35s%s %s\n",
				colorRed, colorReset, colorBold, issue.Key, colorReset, issue.Message)
			errors++
		case validator.SeverityWarning:
			fmt.Printf("%s⚠ WARN%s    %s%-35s%s %s\n",
				colorYellow, colorReset, colorBold, issue.Key, colorReset, issue.Message)
			warnings++
		case validator.SeverityInfo:
			fmt.Printf("%s⚬ INFO%s    %s%-35s%s %s\n",
				colorBlue, colorReset, colorGray, issue.Key, colorReset, issue.Message)
			infos++
		}
	}

	fmt.Println(strings.Repeat("─", 55))

	// TODO: undecided if we want to show 0 counts or just omit them
	if errors > 0 {
		fmt.Printf("  %s%d error(s)%s", colorRed, errors, colorReset)
	}
	if warnings > 0 {
		fmt.Printf("  %s%d warning(s)%s", colorYellow, warnings, colorReset)
	}
	if infos > 0 {
		fmt.Printf("  %s%d info%s", colorBlue, infos, colorReset)
	}
	fmt.Println()

	return errors > 0
}

func PrintJSON(issues []validator.Issue) {
	fmt.Println("[")
	for i, issue := range issues {
		comma := ","
		if i == len(issues)-1 {
			comma = ""
		}
		fmt.Printf("  {\"key\": %q, \"severity\": %q, \"message\": %q}%s\n",
			issue.Key, issue.Severity, issue.Message, comma)
	}
	fmt.Println("]")
}
