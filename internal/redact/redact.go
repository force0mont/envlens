package redact

import (
	"strings"

	"github.com/user/envlens/internal/parser"
)

const redactedValue = "***REDACTED***"

// defaultSecretPatterns are key substrings that indicate sensitive values.
var defaultSecretPatterns = []string{
	"password", "passwd", "secret", "token", "api_key", "apikey",
	"auth", "private_key", "private", "credential", "cert", "jwt",
	"access_key", "signing_key",
}

// Options controls redaction behaviour.
type Options struct {
	// ExtraPatterns are additional key substrings to treat as sensitive.
	ExtraPatterns []string
	// Allowlist contains exact key names that should never be redacted.
	Allowlist []string
}

// Redact returns a copy of env with sensitive values replaced.
func Redact(env parser.EnvFile, opts Options) parser.EnvFile {
	patterns := append(defaultSecretPatterns, opts.ExtraPatterns...)
	allowSet := make(map[string]struct{}, len(opts.Allowlist))
	for _, k := range opts.Allowlist {
		allowSet[k] = struct{}{}
	}

	result := make(parser.EnvFile, len(env))
	for k, v := range env {
		if _, allowed := allowSet[k]; allowed {
			result[k] = v
			continue
		}
		if isSensitive(k, patterns) {
			result[k] = redactedValue
		} else {
			result[k] = v
		}
	}
	return result
}

// isSensitive reports whether key contains any of the given patterns.
func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
