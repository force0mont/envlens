package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/merge"
	"github.com/user/envlens/internal/parser"
)

func init() {
	var labels []string
	var strategy string

	mergeCmd := &cobra.Command{
		Use:   "merge <file1> <file2> [fileN...]",
		Short: "Merge multiple .env files, reporting conflicts",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var envFiles []parser.EnvFile
			for _, path := range args {
				env, err := parser.Parse(path)
				if err != nil {
					return fmt.Errorf("failed to parse %s: %w", path, err)
				}
				envFiles = append(envFiles, env)
			}

			if len(labels) == 0 {
				labels = args
			}

			strat := merge.StrategyFirst
			if strategy == "last" {
				strat = merge.StrategyLast
			}

			result := merge.Merge(envFiles, strat)
			fmt.Print(merge.Format(result, labels))

			if len(result.Conflicts) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	mergeCmd.Flags().StringSliceVarP(&labels, "labels", "l", nil,
		"Comma-separated labels for each input file (defaults to file paths)")
	mergeCmd.Flags().StringVarP(&strategy, "strategy", "s", "first",
		"Conflict resolution strategy: first or last")

	rootCmd.AddCommand(mergeCmd)
}
