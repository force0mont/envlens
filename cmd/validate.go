package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/envlens/envlens/internal/parser"
	"github.com/envlens/envlens/internal/validate"
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate an .env file against a set of rules",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		label, _ := cmd.Flags().GetString("label")
		if label == "" {
			label = filePath
		}
		requiredKeys, _ := cmd.Flags().GetStringSlice("require")
		patternFlags, _ := cmd.Flags().GetStringSlice("pattern")

		env, err := parser.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		var rules []validate.Rule
		for _, key := range requiredKeys {
			rules = append(rules, validate.Rule{Key: strings.TrimSpace(key), Required: true})
		}
		for _, pf := range patternFlags {
			parts := strings.SplitN(pf, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --pattern flag %q: expected KEY=PATTERN", pf)
			}
			rules = append(rules, validate.Rule{
				Key:     strings.TrimSpace(parts[0]),
				Pattern: strings.TrimSpace(parts[1]),
			})
		}

		findings := validate.Validate(env, rules)
		fmt.Print(validate.Format(findings, label))

		if len(findings) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	validateCmd.Flags().String("label", "", "label for the env file (defaults to filename)")
	validateCmd.Flags().StringSlice("require", nil, "comma-separated list of required keys")
	validateCmd.Flags().StringSlice("pattern", nil, "KEY=PATTERN pairs to validate values against")
	rootCmd.AddCommand(validateCmd)
}
