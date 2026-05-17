package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/internal/envclone"
	"github.com/envlens/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	var outputFile string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "envclone <file> <src=dst>...",
		Short: "Clone env keys to new names within a file",
		Long: `Duplicates the value of one or more keys to new key names.
Each pair is specified as SRC=DST. Use --overwrite to replace
an existing destination key.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			pairArgs := args[1:]

			pairs := make(map[string]string, len(pairArgs))
			for _, p := range pairArgs {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return fmt.Errorf("invalid pair %q: expected SRC=DST", p)
				}
				pairs[parts[0]] = parts[1]
			}

			entries, err := parser.Parse(filePath)
			if err != nil {
				return fmt.Errorf("parse %s: %w", filePath, err)
			}

			out, results := envclone.Clone(entries, pairs, overwrite)
			fmt.Print(envclone.Format(results))

			dest := outputFile
			if dest == "" {
				dest = filePath
			}

			var sb strings.Builder
			for _, e := range out {
				sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
			}
			if err := os.WriteFile(dest, []byte(sb.String()), 0644); err != nil {
				return fmt.Errorf("write %s: %w", dest, err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "write result to this file instead of the source")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination key if it already exists")

	rootCmd.AddCommand(cmd)
}
