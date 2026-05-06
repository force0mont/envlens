package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/rename"
)

func init() {
	var pairs []string
	var output string

	cmd := &cobra.Command{
		Use:   "rename <file>",
		Short: "Rename keys in an .env file",
		Long:  "Rename one or more keys in an .env file. Use --pair OLD=NEW flags to specify renames.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse %s: %w", args[0], err)
			}

			pairMap := make(map[string]string, len(pairs))
			for _, p := range pairs {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid pair %q: expected OLD=NEW", p)
				}
				pairMap[parts[0]] = parts[1]
			}

			resulted, results := rename.Rename(env, pairMap)
			fmt.Print(rename.Format(results))

			if output != "" {
				var sb strings.Builder
				for k, v := range resulted {
					fmt.Fprintf(&sb, "%s=%s\n", k, v)
				}
				if err := os.WriteFile(output, []byte(sb.String()), 0644); err != nil {
					return fmt.Errorf("write %s: %w", output, err)
				}
				fmt.Printf("Written to %s\n", output)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&pairs, "pair", "p", nil, "Rename pair in OLD=NEW format (repeatable)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Write resulting env to this file")
	_ = cmd.MarkFlagRequired("pair")

	rootCmd.AddCommand(cmd)
}
