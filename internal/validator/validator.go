package validator

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/norowachi/env-doctor/internal/parser"
)

// Severity levels
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

type Issue struct {
	Key      string
	Severity Severity
	Message  string
}

// Placeholder values that shouldn't be in production
var globalSuspiciousValues = []string{
	"changeme", "secret", "xxx", "todo", "fixme",
	"password", "12345", "test", "example", "replace",
	"your_", "_here", "insert_",
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// compare actual env against the example schema, and return issues
func Validate(example *parser.EnvFile, actual *parser.EnvFile) []Issue {
	var issues []Issue

	for _, key := range example.Order {
		exampleEntry := example.Entries[key]
		actualEntry, exists := actual.Entries[key]

		// Missing key
		if !exists {
			sev := SeverityWarning
			if exampleEntry.Required {
				sev = SeverityError
			}
			issues = append(issues, Issue{
				Key:      key,
				Severity: sev,
				Message:  "Key is missing from your " + actual.Name + " file",
			})
			continue
		}

		// Empty value
		if actualEntry.Value == "" {
			sev := SeverityWarning
			if exampleEntry.Required {
				sev = SeverityError
			}
			issues = append(issues, Issue{
				Key:      key,
				Severity: sev,
				Message:  "Value is empty",
			})
			continue
		}

		// Placeholder value
		lowerVal := strings.ToLower(actualEntry.Value)
		for _, bad := range globalSuspiciousValues {
			if strings.Contains(lowerVal, bad) {
				issues = append(issues, Issue{
					Key:      key,
					Severity: SeverityWarning,
					Message:  "Value looks like a placeholder: \"" + actualEntry.Value + "\"",
				})
				break
			}
		}

		// Type validation
		if exampleEntry.Type != "" {
			if typeIssue := validateType(key, actualEntry.Value, exampleEntry.Type); typeIssue != nil {
				issues = append(issues, *typeIssue)
			}
		}
	}

	// Extra keys
	for _, key := range actual.Order {
		if _, exists := example.Entries[key]; !exists {
			issues = append(issues, Issue{
				Key:      key,
				Severity: SeverityInfo,
				Message:  "Key exists in " + actual.Name + " but not in " + example.Name,
			})
		}
	}

	return issues
}

func validateType(key, value, typeName string) *Issue {
	// Regex type, @type: /pattern/
	if strings.HasPrefix(typeName, "/") && strings.HasSuffix(typeName, "/") {
		pattern := typeName[1 : len(typeName)-1]
		re, err := regexp.Compile(pattern)
		if err != nil {
			return &Issue{Key: key, Severity: SeverityError,
				Message: "Invalid regex pattern in schema: " + pattern + " (" + err.Error() + ")"}
		}
		if !re.MatchString(value) {
			return &Issue{Key: key, Severity: SeverityError,
				Message: "Value \"" + value + "\" does not match pattern /" + pattern + "/"}
		}
		return nil
	}

	switch strings.ToLower(typeName) {
	case "url":
		u, err := url.ParseRequestURI(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return &Issue{Key: key, Severity: SeverityError,
				Message: "Expected a valid URL, got: \"" + value + "\""}
		}

	case "number", "int", "integer":
		if _, err := strconv.Atoi(value); err != nil {
			return &Issue{Key: key, Severity: SeverityError,
				Message: "Expected a number, got: \"" + value + "\""}
		}

	case "boolean", "bool":
		lower := strings.ToLower(value)
		valid := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
		if !valid[lower] {
			return &Issue{Key: key, Severity: SeverityError,
				Message: "Expected a boolean (true/false/1/0), got: \"" + value + "\""}
		}

	case "email":
		if !emailRegex.MatchString(value) {
			return &Issue{Key: key, Severity: SeverityError,
				Message: "Expected a valid email address, got: \"" + value + "\""}
		}

	default:
		return &Issue{Key: key, Severity: SeverityError,
			Message: "Unknown @type \"" + typeName + "\" in schema, valid types: url, number, boolean, email or /regex/"}
	}

	return nil
}
