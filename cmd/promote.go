package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlens/internal/parser"
	"envlens/internal/promote"
)

func init() {
	var srcLabel string
	var dstLabel string
	var keys []string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "promote <src> <dst>",
		Short: "Promote env vars from a source file into a destination file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath, dstPath := args[0], args[1]

			src, err := parser.Parse(srcPath)
			if err != nil {
				return fmt.Errorf("reading source %q: %w", srcPath, err)
			}
			dst, err := parser.Parse(dstPath)
			if err != nil {
				return fmt.Errorf("reading destination %q: %w", dstPath, err)
			}

			if srcLabel == "" {
				srcLabel = srcPath
			}
			if dstLabel == "" {
				dstLabel = dstPath
			}

			merged, results := promote.Promote(src, dst, keys)
			fmt.Print(promote.Format(results, srcLabel, dstLabel))

			if outputFile != "" {
				var sb strings.Builder
				for k, v := range merged {
					sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
				}
				if err := os.WriteFile(outputFile, []byte(sb.String()), 0644); err != nil {
					return fmt.Errorf("writing output file: %w", err)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&srcLabel, "src-label", "", "label for the source env file")
	cmd.Flags().StringVar(&dstLabel, "dst-label", "", "label for the destination env file")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated list of keys to promote (default: all)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "write merged result to this file")

	rootCmd.AddCommand(cmd)
}
