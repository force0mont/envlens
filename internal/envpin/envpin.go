// Package envpin records the current state of an env file and detects
// unexpected changes (drift) against a previously pinned version.
package envpin

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Pin represents a snapshot of key→value pairs taken at a point in time.
type Pin struct {
	CreatedAt time.Time         `json:"created_at"`
	Entries   map[string]string `json:"entries"`
}

// DriftEntry describes a single key that has drifted from its pinned value.
type DriftEntry struct {
	Key    string
	Status string // "changed", "added", "removed"
	Pinned string
	Current string
}

// Create builds a new Pin from the provided env map.
func Create(env map[string]string) Pin {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Pin{CreatedAt: time.Now().UTC(), Entries: copy}
}

// Save writes a Pin to the given file path as JSON.
func Save(p Pin, path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("envpin: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a Pin from the given file path.
func Load(path string) (Pin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Pin{}, fmt.Errorf("envpin: read %s: %w", path, err)
	}
	var p Pin
	if err := json.Unmarshal(data, &p); err != nil {
		return Pin{}, fmt.Errorf("envpin: unmarshal: %w", err)
	}
	return p, nil
}

// Detect compares current env against a pinned state and returns drift entries.
func Detect(pin Pin, current map[string]string) []DriftEntry {
	var results []DriftEntry

	for k, pinned := range pin.Entries {
		cur, ok := current[k]
		if !ok {
			results = append(results, DriftEntry{Key: k, Status: "removed", Pinned: pinned})
		} else if cur != pinned {
			results = append(results, DriftEntry{Key: k, Status: "changed", Pinned: pinned, Current: cur})
		}
	}
	for k, cur := range current {
		if _, ok := pin.Entries[k]; !ok {
			results = append(results, DriftEntry{Key: k, Status: "added", Current: cur})
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })
	return results
}

// Format returns a human-readable drift report.
func Format(entries []DriftEntry) string {
	if len(entries) == 0 {
		return "no drift detected\n"
	}
	out := fmt.Sprintf("%d drift(s) detected:\n", len(entries))
	for _, e := range entries {
		switch e.Status {
		case "added":
			out += fmt.Sprintf("  + %-30s (new value: %q)\n", e.Key, e.Current)
		case "removed":
			out += fmt.Sprintf("  - %-30s (was: %q)\n", e.Key, e.Pinned)
		case "changed":
			out += fmt.Sprintf("  ~ %-30s %q → %q\n", e.Key, e.Pinned, e.Current)
		}
	}
	return out
}
