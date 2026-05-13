package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/envlens/internal/parser"
	"github.com/envlens/envlens/internal/required"
	"github.com/spf13/cobra"
)

func init() {
	var keys []string

	cmd := &cobra.Command{
		Use:   "required <file>",
		Short: "Check that required keys exist and are non-empty in an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(keys) == 0 {
				return fmt.Errorf("at least one --key must be specified")
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}
			defer f.Close()

			entries, err := parser.Parse(f)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			// Support comma-separated keys within a single --key flag value.
			var expanded []string
			for _, k := range keys {
				for _, part := range strings.Split(k, ",") {
					if p := strings.TrimSpace(part); p != "" {
						expanded = append(expanded, p)
					}
				}
			}

			result := required.Check(entries, expanded)
			fmt.Print(required.Format(result))

			if len(result.Findings) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&keys, "key", "k", nil,
		"required key name (repeatable; comma-separated values accepted)")

	rootCmd.AddCommand(cmd)
}
