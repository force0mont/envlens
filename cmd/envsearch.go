package cmd

import (
	"fmt"
	"os"

	"github.com/envlens/internal/envsearch"
	"github.com/envlens/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	var field string
	var caseSensitive bool

	cmd := &cobra.Command{
		Use:   "search <pattern> <file>",
		Short: "Search for keys or values matching a pattern in an env file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]
			filePath := args[1]

			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("cannot open file %q: %w", filePath, err)
			}
			defer f.Close()

			entries, err := parser.Parse(f)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			var mf envsearch.MatchField
			switch field {
			case "key":
				mf = envsearch.MatchKey
			case "value":
				mf = envsearch.MatchValue
			default:
				mf = envsearch.MatchBoth
			}

			results, err := envsearch.Search(entries, envsearch.Options{
				Pattern:       pattern,
				Field:         mf,
				CaseSensitive: caseSensitive,
			})
			if err != nil {
				return err
			}

			fmt.Print(envsearch.Format(results, pattern))
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "both", "Field to search: key, value, or both")
	cmd.Flags().BoolVar(&caseSensitive, "case-sensitive", false, "Enable case-sensitive matching")
	rootCmd.AddCommand(cmd)
}
