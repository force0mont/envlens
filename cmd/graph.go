package cmd

import (
	"fmt"
	"os"

	"github.com/envlens/internal/graph"
	"github.com/envlens/internal/parser"
	"github.com/spf13/cobra"
)

var graphLabels []string

var graphCmd = &cobra.Command{
	Use:   "graph [file1] [file2] ...",
	Short: "Show a relationship graph of keys across multiple .env files",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(graphLabels) > 0 && len(graphLabels) != len(args) {
			return fmt.Errorf("number of --labels must match number of files")
		}

		files := make(map[string]parser.EnvFile, len(args))
		for i, path := range args {
			env, err := parser.Parse(path)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
			label := path
			if i < len(graphLabels) {
				label = graphLabels[i]
			}
			files[label] = env
		}

		g := graph.Build(files)
		fmt.Print(graph.Format(g))
		return nil
	},
}

func init() {
	graphCmd.Flags().StringSliceVar(&graphLabels, "labels", nil, "comma-separated labels for each file")
	if err := rootCmd.AddCommand(graphCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
