package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envlens/internal/envdoc"
)

func init() {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "envdoc <file>",
		Short: "Generate documentation from a .env file's inline comments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			lines, err := readLines(path)
			if err != nil {
				return err
			}

			doc := envdoc.Generate(lines)

			switch outputFormat {
			case "markdown", "":
				fmt.Print(envdoc.Format(doc))
			case "count":
				fmt.Printf("%d variable(s) documented\n", len(doc.Entries))
			default:
				return fmt.Errorf("unknown format %q (supported: markdown, count)", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "markdown", "Output format: markdown, count")
	rootCmd.AddCommand(cmd)
}

// readLines opens the file at path and returns its contents as a slice of strings,
// one entry per line. It returns an error if the file cannot be opened or read.
func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return lines, nil
}
