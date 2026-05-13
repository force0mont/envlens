package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/envlens/internal/groupby"
	"github.com/envlens/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	var sep string
	var depth int

	cmd := &cobra.Command{
		Use:   "groupby <file>",
		Short: "Group environment variables by key prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse %q: %w", args[0], err)
			}

			result := groupby.ByPrefix(entries, sep, depth)
			out := groupby.Format(result)

			fmt.Print(out)

			total := 0
			for _, g := range result.Groups {
				total += len(g.Entries)
			}
			total += len(result.Ungrouped)

			fmt.Fprintf(os.Stderr, "%d group(s), %d ungrouped, %d total\n",
				len(result.Groups), len(result.Ungrouped), total)

			return nil
		},
	}

	cmd.Flags().StringVarP(&sep, "sep", "s", "_", "separator used to split key into prefix segments")
	cmd.Flags().IntVarP(&depth, "depth", "d", 1, "number of prefix segments to use for grouping")

	// allow depth to be passed as positional-style override via env for scripting
	if v := os.Getenv("ENVLENS_GROUPBY_DEPTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			depth = n
		}
	}

	rootCmd.AddCommand(cmd)
}
