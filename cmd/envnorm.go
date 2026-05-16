package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/envnorm"
	"github.com/user/envlens/internal/parser"
)

func init() {
	var rules []string
	var output string

	cmd := &cobra.Command{
		Use:   "envnorm <file>",
		Short: "Normalise keys and values in a .env file",
		Long: `Apply normalisation rules to keys and values.

Available rules: uppercase, lowercase, trim_values, snake_case`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse %s: %w", args[0], err)
			}

			parsedRules := make([]envnorm.Rule, 0, len(rules))
			for _, r := range rules {
				switch r {
				case "uppercase", "lowercase", "trim_values", "snake_case":
					parsedRules = append(parsedRules, envnorm.Rule(r))
				default:
					return fmt.Errorf("unknown rule %q", r)
				}
			}

			if len(parsedRules) == 0 {
				return fmt.Errorf("at least one --rule is required")
			}

			results := envnorm.Normalise(entries, parsedRules)
			fmt.Print(envnorm.Format(results))

			if output != "" {
				normalised := envnorm.ToEntries(results)
				f, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
				defer f.Close()
				for _, e := range normalised {
					fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&rules, "rule", nil, "normalisation rule(s) to apply (repeatable)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write normalised env to this file")
	rootCmd.AddCommand(cmd)
}
