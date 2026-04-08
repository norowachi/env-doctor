package validator_test

import (
	"os"
	"testing"

	"github.com/norowachi/env-doctor/internal/parser"
	"github.com/norowachi/env-doctor/internal/validator"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "env-doctor-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestMissingKey(t *testing.T) {
	exPath := writeTemp(t, "# @required\nDB_URL=\n")
	envPath := writeTemp(t, "OTHER=value\n")
	defer os.Remove(exPath)
	defer os.Remove(envPath)

	example, _ := parser.Parse(exPath)
	actual, _ := parser.Parse(envPath)
	issues := validator.Validate(example, actual)

	if len(issues) == 0 {
		t.Fatal("expected issues, got none")
	}
	found := false
	for _, i := range issues {
		if i.Key == "DB_URL" && i.Severity == validator.SeverityError {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ERROR for missing DB_URL, got: %+v", issues)
	}
}

func TestTypeURL(t *testing.T) {
	exPath := writeTemp(t, "# @type: url\nAPI_URL=\n")
	envPath := writeTemp(t, "API_URL=not-a-url\n")
	defer os.Remove(exPath)
	defer os.Remove(envPath)

	example, _ := parser.Parse(exPath)
	actual, _ := parser.Parse(envPath)
	issues := validator.Validate(example, actual)

	found := false
	for _, i := range issues {
		if i.Key == "API_URL" && i.Severity == validator.SeverityError {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ERROR for bad URL, got: %+v", issues)
	}
}

func TestValidURL(t *testing.T) {
	exPath := writeTemp(t, "# @type: url\nAPI_URL=\n")
	envPath := writeTemp(t, "API_URL=https://api.example.com\n")
	defer os.Remove(exPath)
	defer os.Remove(envPath)

	example, _ := parser.Parse(exPath)
	actual, _ := parser.Parse(envPath)
	issues := validator.Validate(example, actual)

	for _, i := range issues {
		if i.Key == "API_URL" && i.Severity == validator.SeverityError {
			t.Errorf("unexpected ERROR for valid URL: %+v", i)
		}
	}
}

func TestSuspiciousValue(t *testing.T) {
	exPath := writeTemp(t, "APP_SECRET=\n")
	envPath := writeTemp(t, "APP_SECRET=changeme\n")
	defer os.Remove(exPath)
	defer os.Remove(envPath)

	example, _ := parser.Parse(exPath)
	actual, _ := parser.Parse(envPath)
	issues := validator.Validate(example, actual)

	found := false
	for _, i := range issues {
		if i.Key == "APP_SECRET" && i.Severity == validator.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Errorf("expected WARNING for suspicious value, got: %+v", issues)
	}
}

func TestAllGood(t *testing.T) {
	exPath := writeTemp(t, "# @type: number\nPORT=\n# @type: boolean\nDEBUG=\n")
	envPath := writeTemp(t, "PORT=3000\nDEBUG=true\n")
	defer os.Remove(exPath)
	defer os.Remove(envPath)

	example, _ := parser.Parse(exPath)
	actual, _ := parser.Parse(envPath)
	issues := validator.Validate(example, actual)

	for _, i := range issues {
		if i.Severity == validator.SeverityError || i.Severity == validator.SeverityWarning {
			t.Errorf("unexpected issue: %+v", i)
		}
	}
}
