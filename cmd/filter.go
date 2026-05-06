package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/envlens/internal/filter"
	"github.com/yourusername/envlens/internal/parser"
)

func init() {
	var prefix string
	var pattern string
	var keys string

	cmd := &cobra.Command{
		Use:   "filter <file>",
		Short: "Filter environment variables by prefix, pattern, or key list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					k = strings.TrimSpace(k)
					if k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			result, err := filter.Filter(env, filter.Options{
				Prefix:  prefix,
				Pattern: pattern,
				Keys:    keyList,
			})
			if err != nil {
				return err
			}

			fmt.Fprint(os.Stdout, filter.Format(result))
			return nil
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Filter keys by prefix")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Filter keys by regex pattern")
	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of exact keys to include")

	rootCmd.AddCommand(cmd)
}
