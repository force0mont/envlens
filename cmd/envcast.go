package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/envlens/internal/envcast"
	"github.com/envlens/envlens/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	var rules []string

	cmd := &cobra.Command{
		Use:   "envcast <file>",
		Short: "Cast env variable values to specified types and report results",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			parsedRules, err := parseRules(rules)
			if err != nil {
				return err
			}

			results := envcast.Cast(entries, parsedRules)
			fmt.Print(envcast.Format(results))

			for _, r := range results {
				if r.Err != nil {
					os.Exit(1)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&rules, "rule", "r", nil,
		"Cast rule in KEY:TYPE format (e.g. PORT:int). Repeat for multiple rules.")
	_ = cmd.MarkFlagRequired("rule")

	rootCmd.AddCommand(cmd)
}

// parseRules converts "KEY:TYPE" strings into envcast.Rule values.
func parseRules(raw []string) ([]envcast.Rule, error) {
	rules := make([]envcast.Rule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid rule %q: expected KEY:TYPE", s)
		}
		t := envcast.TargetType(strings.ToLower(parts[1]))
		switch t {
		case envcast.TypeString, envcast.TypeInt, envcast.TypeFloat, envcast.TypeBool:
		default:
			return nil, fmt.Errorf("unknown type %q in rule %q", parts[1], s)
		}
		rules = append(rules, envcast.Rule{Key: parts[0], Type: t})
	}
	return rules, nil
}
