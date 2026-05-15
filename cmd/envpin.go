package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envlens/internal/envpin"
	"github.com/yourorg/envlens/internal/parser"
)

func init() {
	var pinFile string

	cmd := &cobra.Command{
		Use:   "envpin",
		Short: "Pin or check drift of environment variables",
	}

	saveCmd := &cobra.Command{
		Use:   "save <env-file>",
		Short: "Save a pin of the current env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			env := make(map[string]string, len(entries))
			for _, e := range entries {
				env[e.Key] = e.Value
			}
			p := envpin.Create(env)
			if err := envpin.Save(p, pinFile); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "pin saved to %s\n", pinFile)
			return nil
		},
	}
	saveCmd.Flags().StringVarP(&pinFile, "output", "o", ".envpin.json", "path to write the pin file")

	checkCmd := &cobra.Command{
		Use:   "check <env-file>",
		Short: "Detect drift between current env file and a saved pin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			env := make(map[string]string, len(entries))
			for _, e := range entries {
				env[e.Key] = e.Value
			}
			pin, err := envpin.Load(pinFile)
			if err != nil {
				return err
			}
			drift := envpin.Detect(pin, env)
			fmt.Fprint(cmd.OutOrStdout(), envpin.Format(drift))
			if len(drift) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}
	checkCmd.Flags().StringVarP(&pinFile, "pin", "p", ".envpin.json", "path to the pin file")

	cmd.AddCommand(saveCmd, checkCmd)
	rootCmd.AddCommand(cmd)
}
