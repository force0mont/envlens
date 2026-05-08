package setop

import "sort"

// EnvFile represents a parsed environment file as a key-value map.
type EnvFile = map[string]string

// Result holds the outcome of a set operation.
type Result struct {
	Keys   []string
	Values map[string]string
}

// Intersect returns keys (and their values from the first file) present in ALL provided env maps.
func Intersect(files ...EnvFile) Result {
	if len(files) == 0 {
		return Result{Values: map[string]string{}}
	}
	counts := make(map[string]int)
	for _, f := range files {
		for k := range f {
			counts[k]++
		}
	}
	result := Result{Values: map[string]string{}}
	for k, c := range counts {
		if c == len(files) {
			result.Keys = append(result.Keys, k)
			result.Values[k] = files[0][k]
		}
	}
	sort.Strings(result.Keys)
	return result
}

// Union returns all keys present in ANY of the provided env maps.
// Values are taken from the first file that defines the key.
func Union(files ...EnvFile) Result {
	result := Result{Values: map[string]string{}}
	seen := make(map[string]bool)
	for _, f := range files {
		for k, v := range f {
			if !seen[k] {
				seen[k] = true
				result.Keys = append(result.Keys, k)
				result.Values[k] = v
			}
		}
	}
	sort.Strings(result.Keys)
	return result
}

// Difference returns keys present in the first file but NOT in any of the others.
func Difference(base EnvFile, others ...EnvFile) Result {
	result := Result{Values: map[string]string{}}
	for k, v := range base {
		found := false
		for _, o := range others {
			if _, ok := o[k]; ok {
				found = true
				break
			}
		}
		if !found {
			result.Keys = append(result.Keys, k)
			result.Values[k] = v
		}
	}
	sort.Strings(result.Keys)
	return result
}

// Format renders a Result as KEY=VALUE lines.
func Format(r Result) string {
	if len(r.Keys) == 0 {
		return "(no keys)\n"
	}
	out := ""
	for _, k := range r.Keys {
		out += k + "=" + r.Values[k] + "\n"
	}
	return out
}
