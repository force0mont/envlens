package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envlens/envlens/internal/diff"
	"github.com/envlens/envlens/internal/parser"
)

var diffCmd = &cobra.Command{
	Use:   "diff <file1> <file2>",
	Short: "Diff two .env files and show added, removed, and changed keys",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		file1, err := parser.Parse(args[0])
		if err != nil {
			return fmt.Errorf("reading %s: %w", args[0], err)
		}

		file2, err := parser.Parse(args[1])
		if err != nil {
			return fmt.Errorf("reading %s: %w", args[1], err)
		}

		result := diff.Diff(file1, file2)

		if len(result) == 0 {
			fmt.Println("No differences found.")
			return nil
		}

		labelA, _ := cmd.Flags().GetString("label-a")
		labelB, _ := cmd.Flags().GetString("label-b")

		if labelA == "" {
			labelA = args[0]
		}
		if labelB == "" {
			labelB = args[1]
		}

		output := diff.Format(result, labelA, labelB)
		fmt.Print(output)

		exitCode, _ := cmd.Flags().GetBool("exit-code")
		if exitCode {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	diffCmd.Flags().String("label-a", "", "Label for the first file (default: file path)")
	diffCmd.Flags().String("label-b", "", "Label for the second file (default: file path)")
	diffCmd.Flags().Bool("exit-code", false, "Exit with code 1 if differences are found")
	rootCmd.AddCommand(diffCmd)
}
