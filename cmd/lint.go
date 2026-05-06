package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envlens/internal/lint"
	"github.com/yourorg/envlens/internal/parser"
)

func init() {
	var strict bool

	lintCmd := &cobra.Command{
		Use:   "lint <file>",
		Short: "Lint an .env file for style and correctness issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			env, err := parser.Parse(path)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}

			findings := lint.Lint(env)
			fmt.Print(lint.Format(findings))

			if strict && len(findings) > 0 {
				os.Exit(1)
			}

			return nil
		},
	}

	lintCmd.Flags().BoolVar(&strict, "strict", false, "exit with code 1 if any issues are found")

	rootCmd.AddCommand(lintCmd)
}
