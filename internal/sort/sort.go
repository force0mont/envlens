package sort

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Order defines the sort direction.
type Order string

const (
	Ascending  Order = "asc"
	Descending Order = "desc"
)

// Result holds a sorted list of env entries.
type Result struct {
	Entries []parser.Entry
	Order   Order
}

// Sort returns the env entries sorted by key in the specified order.
func Sort(env []parser.Entry, order Order) Result {
	copied := make([]parser.Entry, len(env))
	copy(copied, env)

	sort.Slice(copied, func(i, j int) bool {
		a := strings.ToLower(copied[i].Key)
		b := strings.ToLower(copied[j].Key)
		if order == Descending {
			return a > b
		}
		return a < b
	})

	return Result{Entries: copied, Order: order}
}

// Format renders the sorted entries as KEY=VALUE lines.
func Format(r Result) string {
	if len(r.Entries) == 0 {
		return "(no entries)\n"
	}
	var sb strings.Builder
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}
