package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/compare"
	"github.com/user/envlens/internal/parser"
)

func init() {
	var labels []string

	cmd := &cobra.Command{
		Use:   "compare <file1> <file2> [fileN...]",
		Short: "Compare key values across multiple .env files",
		Long: `Compare reads every key present in any of the supplied .env files
and shows whether its value is consistent across all files.
Keys that differ between files are highlighted with a '!' marker.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(labels) > 0 && len(labels) != len(args) {
				return fmt.Errorf("number of labels (%d) must match number of files (%d)", len(labels), len(args))
			}

			envs := make(map[string]parser.EnvFile, len(args))
			effectiveLabels := make([]string, len(args))

			for i, path := range args {
				lbl := path
				if i < len(labels) {
					lbl = labels[i]
				}
				effectiveLabels[i] = lbl

				env, err := parser.Parse(path)
				if err != nil {
					return fmt.Errorf("parsing %s: %w", path, err)
				}
				envs[lbl] = env
			}

			results := compare.Compare(envs)
			fmt.Fprint(os.Stdout, compare.Format(results, effectiveLabels))
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&labels, "labels", "l", nil, "comma-separated labels for each file (must match file count)")
	rootCmd.AddCommand(cmd)
}
