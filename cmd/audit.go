package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/audit"
	"github.com/user/envlens/internal/parser"
)

var auditCmd = &cobra.Command{
	Use:   "audit [file...]",
	Short: "Audit one or more .env files for common issues",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		exitCode := 0

		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("cannot open %q: %w", path, err)
			}
			env, err := parser.Parse(f)
			f.Close()
			if err != nil {
				return fmt.Errorf("cannot parse %q: %w", path, err)
			}

			report := audit.Audit(path, env)
			fmt.Print(audit.Format(report))

			if report.HasIssues() {
				exitCode = 1
			}
		}

		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}
