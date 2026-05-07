package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envlens/internal/copy"
	"github.com/yourorg/envlens/internal/parser"
)

func init() {
	var srcLabel string
	var dstLabel string
	var overwrite bool
	var outputFile string

	cmd := &cobra.Command{
		Use:   "copy <src.env> <dst.env> KEY[,KEY...]",
		Short: "Copy specific keys from one .env file into another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath, dstPath, rawKeys := args[0], args[1], args[2]

			src, err := parser.Parse(srcPath)
			if err != nil {
				return fmt.Errorf("reading source: %w", err)
			}
			dst, err := parser.Parse(dstPath)
			if err != nil {
				return fmt.Errorf("reading destination: %w", err)
			}

			keys := strings.Split(rawKeys, ",")
			for i, k := range keys {
				keys[i] = strings.TrimSpace(k)
			}

			if srcLabel == "" {
				srcLabel = srcPath
			}
			if dstLabel == "" {
				dstLabel = dstPath
			}

			out, results := copy.Copy(src, dst, keys, overwrite)
			fmt.Print(copy.Format(results, srcLabel, dstLabel))

			if outputFile != "" {
				var sb strings.Builder
				for k, v := range out {
					sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
				}
				if err := os.WriteFile(outputFile, []byte(sb.String()), 0644); err != nil {
					return fmt.Errorf("writing output: %w", err)
				}
				fmt.Printf("\nwrote merged env to %s\n", outputFile)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&srcLabel, "src-label", "", "label for source file")
	cmd.Flags().StringVar(&dstLabel, "dst-label", "", "label for destination file")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in destination")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "write result to file instead of stdout")

	rootCmd.AddCommand(cmd)
}
