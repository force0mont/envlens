package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourorg/envlens/internal/mask"
	"github.com/yourorg/envlens/internal/parser"
)

func init() {
	var extraPatterns []string
	var visibleChars int

	cmd := &cobra.Command{
		Use:   "mask <file>",
		Short: "Print env file with sensitive values partially masked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}

			r := mask.Mask(env, extraPatterns, visibleChars)
			fmt.Print(mask.Format(r))

			fmt.Fprintf(os.Stderr, "# %s masked\n", pluralise(r.MaskCount, "key"))
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&extraPatterns, "pattern", "p", nil,
		"additional key substrings to treat as sensitive (repeatable)")
	cmd.Flags().IntVarP(&visibleChars, "visible", "v", 2,
		"number of leading value characters to leave unmasked")

	rootCmd.AddCommand(cmd)
}

func pluralise(n int, word string) string {
	if n == 1 {
		return strconv.Itoa(n) + " " + word
	}
	return strconv.Itoa(n) + " " + word + "s"
}
