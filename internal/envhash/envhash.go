package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Result holds the hash output for a single env file.
type Result struct {
	Label   string
	Hash    string
	KeyCount int
}

// Hash computes a deterministic SHA-256 digest over the key=value pairs
// in the provided env map. Keys are sorted before hashing so that
// insertion order does not affect the digest.
func Hash(label string, env []parser.Entry) Result {
	keys := make([]string, 0, len(env))
	lookup := make(map[string]string, len(env))
	for _, e := range env {
		keys = append(keys, e.Key)
		lookup[e.Key] = e.Value
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, lookup[k])
	}

	return Result{
		Label:    label,
		Hash:     hex.EncodeToString(h.Sum(nil)),
		KeyCount: len(keys),
	}
}

// Format renders a human-readable table of hash results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no files to hash\n"
	}

	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "%-30s  %s  (%d keys)\n", r.Label, r.Hash, r.KeyCount)
	}

	// Append a note when all hashes are identical.
	if allSame(results) {
		sb.WriteString("\n✓ all files produce the same hash\n")
	} else {
		sb.WriteString("\n✗ files differ\n")
	}
	return sb.String()
}

func allSame(results []Result) bool {
	if len(results) == 0 {
		return true
	}
	first := results[0].Hash
	for _, r := range results[1:] {
		if r.Hash != first {
			return false
		}
	}
	return true
}
