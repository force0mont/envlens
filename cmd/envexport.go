package cmd

import (
	"fmt"
	"os"

	"github.com/envlens/envlens/internal/envexport"
	"github.com/envlens/envlens/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	var format string
	var prefix string
	var omitEmpty bool
	var outputFile string

	cmd := &cobra.Command{
		Use:   "envexport <file>",
		Short: "Export env variables in dotenv, shell, docker, or JSON format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse %q: %w", args[0], err)
			}

			res, err := envexport.Export(entries, envexport.Options{
				OutputFormat: envexport.Format(format),
				Prefix:       prefix,
				OmitEmpty:    omitEmpty,
			})
			if err != nil {
				return err
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(res.Content), 0o644); err != nil {
					return fmt.Errorf("write output: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exported %d variable(s) to %s\n", res.Count, outputFile)
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), res.Content)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv, shell, docker, json")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Prefix to prepend to all keys")
	cmd.Flags().BoolVar(&omitEmpty, "omit-empty", false, "Omit variables with empty values")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")

	rootCmd.AddCommand(cmd)
}
