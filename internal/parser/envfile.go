package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvEntry represents a single key-value pair from an env file.
type EnvEntry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile represents a parsed .env file.
type EnvFile struct {
	Path    string
	Entries []EnvEntry
	Index   map[string]EnvEntry
}

// Parse reads and parses a .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	ef := &EnvFile{
		Path:  path,
		Index: make(map[string]EnvEntry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entry, err := parseLine(line, lineNum)
		if err != nil {
			return nil, fmt.Errorf("parse error at %s:%d: %w", path, lineNum, err)
		}

		ef.Entries = append(ef.Entries, entry)
		ef.Index[entry.Key] = entry
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file %q: %w", path, err)
	}

	return ef, nil
}

func parseLine(line string, lineNum int) (EnvEntry, error) {
	// Strip inline comment
	comment := ""
	if idx := strings.Index(line, " #"); idx != -1 {
		comment = strings.TrimSpace(line[idx+2:])
		line = strings.TrimSpace(line[:idx])
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return EnvEntry{}, fmt.Errorf("invalid format: expected KEY=VALUE")
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

	if key == "" {
		return EnvEntry{}, fmt.Errorf("empty key")
	}

	return EnvEntry{Key: key, Value: value, Comment: comment, Line: lineNum}, nil
}
