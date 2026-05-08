package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/envlens/internal/parser"
	"github.com/yourusername/envlens/internal/setop"
)

func init() {
	var operation string

	cmd := &cobra.Command{
		Use:   "setop <file1> <file2> [fileN...]",
		Short: "Perform set operations (intersect, union, difference) across env files",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			files := make([]setop.EnvFile, 0, len(args))
			for _, path := range args {
				entries, err := parser.Parse(path)
				if err != nil {
					return fmt.Errorf("reading %s: %w", path, err)
				}
				m := make(setop.EnvFile, len(entries))
				for _, e := range entries {
					m[e.Key] = e.Value
				}
				files = append(files, m)
			}

			var result setop.Result
			switch strings.ToLower(operation) {
			case "intersect":
				result = setop.Intersect(files...)
			case "union":
				result = setop.Union(files...)
			case "difference":
				result = setop.Difference(files[0], files[1:]...)
			default:
				return fmt.Errorf("unknown operation %q: choose intersect, union, or difference", operation)
			}

			fmt.Fprint(os.Stdout, setop.Format(result))
			return nil
		},
	}

	cmd.Flags().StringVarP(&operation, "op", "o", "intersect",
		"Set operation to perform: intersect, union, difference")

	rootCmd.AddCommand(cmd)
}
