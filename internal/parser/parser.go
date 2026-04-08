package parser

import (
	"bufio"
	"os"
	"strings"
)

// Properties of each env variable entry from the annotations above it
type Entry struct {
	Key      string
	Value    string
	Required bool
	// string, url, number, boolean, email
	// TODO: support custom types with regex validation
	Type    string
	LineNum int
}

// Parsed entries from env file by order
type EnvFile struct {
	Entries map[string]*Entry
	Order   []string
}

// Read target env and schema files and extracts entries & properties
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ef := &EnvFile{
		Entries: make(map[string]*Entry),
	}

	var pendingEntry Entry
	lineNum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines
		if line == "" {
			pendingEntry = Entry{}
			continue
		}

		// Annotation comment treatment
		if after, ok := strings.CutPrefix(line, "#"); ok {
			annotation := strings.TrimSpace(after)

			if after, ok := strings.CutPrefix(annotation, "@type:"); ok {
				pendingEntry.Type = strings.TrimSpace(after)
			} else if strings.HasPrefix(annotation, "@required") {
				pendingEntry.Required = true
			}
			continue
		}

		// Key=Value line
		if before, after, ok := strings.Cut(line, "="); ok {
			key := strings.TrimSpace(before)
			value := strings.TrimSpace(after)

			// Strip inline quotes
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}

			entry := &Entry{
				Key:      key,
				Value:    value,
				Required: pendingEntry.Required,
				Type:     pendingEntry.Type,
				LineNum:  lineNum,
			}

			ef.Entries[key] = entry
			ef.Order = append(ef.Order, key)

			// Reset
			pendingEntry = Entry{}
		}
	}

	return ef, scanner.Err()
}
