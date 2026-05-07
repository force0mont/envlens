package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/convert"
	"github.com/user/envlens/internal/parser"
)

func init() {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "convert <file>",
		Short: "Convert a .env file to another format (env, export, json, yaml)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			env, err := parser.Parse(path)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}

			fmt := convert.Format(strings.ToLower(outputFormat))
			result, err := convert.Convert(env, fmt)
			if err != nil {
				return err
			}

			_, err = os.Stdout.WriteString(result.Output)
			return err
		},
	}

	cmd.Flags().StringVarP(
		&outputFormat,
		"format", "f",
		"env",
		"Output format: env, export, json, yaml",
	)

	rootCmd.AddCommand(cmd)
}
