package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlens/internal/parser"
	"envlens/internal/scope"
)

func init() {
	var labels []string

	cmd := &cobra.Command{
		Use:   "scope <file1> [file2 ...]",
		Short: "Tag each key with its originating scope (last scope wins on collision)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(labels) == 0 {
				for _, a := range args {
					labels = append(labels, a)
				}
			}
			if len(labels) != len(args) {
				return fmt.Errorf("number of --labels must match number of files")
			}

			envs := make([]map[string]string, 0, len(args))
			for _, path := range args {
				entries, err := parser.Parse(path)
				if err != nil {
					return fmt.Errorf("parsing %s: %w", path, err)
				}
				m := make(map[string]string, len(entries))
				for _, e := range entries {
					m[e.Key] = e.Value
				}
				envs = append(envs, m)
			}

			r, err := scope.Tag(envs, labels)
			if err != nil {
				return err
			}
			fmt.Fprint(os.Stdout, scope.Format(r))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&labels, "labels", nil,
		"comma-separated scope labels (defaults to file paths)")

	// strip brackets/spaces that StringSliceVar adds when printing
	_ = strings.TrimSpace

	rootCmd.AddCommand(cmd)
}
