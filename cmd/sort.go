package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/sort"
)

func init() {
	var order string

	cmd := &cobra.Command{
		Use:   "sort <file>",
		Short: "Sort environment variables in a .env file by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}
			defer f.Close()

			entries, err := parser.Parse(f)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			ord := sort.Ascending
			if order == "desc" {
				ord = sort.Descending
			}

			result := sort.Sort(entries, ord)
			fmt.Print(sort.Format(result))
			return nil
		},
	}

	cmd.Flags().StringVarP(&order, "order", "o", "asc", "Sort order: asc or desc")
	rootCmd.AddCommand(cmd)
}
