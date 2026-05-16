package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/internal/envreplace"
	"github.com/envlens/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	var outputFile string
	var keys []string
	var literal bool

	cmd := &cobra.Command{
		Use:   "envreplace <file> <old> <new>",
		Short: "Find and replace values in a .env file",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, old, newVal := args[0], args[1], args[2]

			entries, err := parser.Parse(filePath)
			if err != nil {
				return fmt.Errorf("parse %q: %w", filePath, err)
			}

			result := envreplace.Replace(entries, old, newVal, keys, literal)

			if outputFile != "" {
				if err := writeEnvFile(outputFile, result.Entries); err != nil {
					return fmt.Errorf("write output: %w", err)
				}
				fmt.Fprint(cmd.OutOrStdout(), envreplace.Format(result))
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), envreplace.Format(result))
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write updated entries to file")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "Restrict replacement to these keys (comma-separated)")
	cmd.Flags().BoolVar(&literal, "literal", false, "Match entire value exactly instead of substring")

	rootCmd.AddCommand(cmd)
}

// writeEnvFile serialises entries back to KEY=VALUE format.
func writeEnvFile(path string, entries []parser.Entry) error {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
