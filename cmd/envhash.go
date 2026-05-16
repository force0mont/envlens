package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/envhash"
	"github.com/user/envlens/internal/parser"
)

func init() {
	var labels []string

	cmd := &cobra.Command{
		Use:   "envhash <file> [file...]",
		Short: "Compute a deterministic SHA-256 hash for one or more .env files",
		Long: `envhash reads each .env file, sorts its key=value pairs, and
produces a SHA-256 digest. Identical digests confirm two files are
functionally equivalent regardless of key ordering or blank lines.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Pad or generate labels.
			for len(labels) < len(args) {
				labels = append(labels, args[len(labels)])
			}

			results := make([]envhash.Result, 0, len(args))
			for i, path := range args {
				f, err := os.Open(path)
				if err != nil {
					return fmt.Errorf("open %s: %w", path, err)
				}
				entries, err := parser.Parse(f)
				f.Close()
				if err != nil {
					return fmt.Errorf("parse %s: %w", path, err)
				}
				results = append(results, envhash.Hash(labels[i], entries))
			}

			fmt.Print(envhash.Format(results))
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil,
		"comma-separated labels for each file (defaults to file paths)")

	rootCmd.AddCommand(cmd)
}
