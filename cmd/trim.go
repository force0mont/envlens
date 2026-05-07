package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envlens/internal/parser"
	"envlens/internal/trim"
)

func init() {
	var outputFile string
	var inPlace bool

	cmd := &cobra.Command{
		Use:   "trim <file>",
		Short: "Remove leading/trailing whitespace from all env values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}

			trimmed, results := trim.Trim(env)

			fmt.Print(trim.Format(results))

			dest := outputFile
			if inPlace {
				dest = args[0]
			}

			if dest == "" {
				return nil
			}

			f, err := os.Create(dest)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer f.Close()

			for k, v := range trimmed {
				if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
					return fmt.Errorf("writing output: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write trimmed values to this file")
	cmd.Flags().BoolVar(&inPlace, "in-place", false, "Overwrite the source file with trimmed values")

	rootCmd.AddCommand(cmd)
}
