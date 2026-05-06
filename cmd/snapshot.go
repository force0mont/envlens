package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/snapshot"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Save or load snapshots of .env files",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save <envfile> <output.json>",
	Short: "Save a .env file as a named snapshot",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		label, _ := cmd.Flags().GetString("label")
		env, err := parser.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		if label == "" {
			label = args[0]
		}
		s := snapshot.New(label, env)
		if err := snapshot.Save(s, args[1]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Snapshot saved to %s (label: %s, %d keys)\n", args[1], s.Label, len(s.Env))
		return nil
	},
}

var snapshotInfoCmd = &cobra.Command{
	Use:   "info <snapshot.json>",
	Short: "Print metadata and keys from a snapshot file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := snapshot.Load(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Label:   %s\n", s.Label)
		fmt.Fprintf(os.Stdout, "Created: %s\n", s.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
		fmt.Fprintf(os.Stdout, "Keys:    %d\n", len(s.Env))
		for k := range s.Env {
			fmt.Fprintf(os.Stdout, "  %s\n", k)
		}
		return nil
	},
}

func init() {
	snapshotSaveCmd.Flags().String("label", "", "Label for the snapshot (defaults to filename)")
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotInfoCmd)
	rootCmd.AddCommand(snapshotCmd)
}
