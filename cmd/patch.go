package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envlens/internal/parser"
	"envlens/internal/patch"
)

func init() {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "patch <file> [ops...]",
		Short: "Apply set/delete/rename operations to an env file",
		Long: `Apply patch operations to an env file.

Operations:
  set:KEY=VALUE    — set or overwrite a key
  delete:KEY       — remove a key
  rename:OLD=NEW   — rename a key`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			opsArgs := args[1:]

			env, err := parser.Parse(filePath)
			if err != nil {
				return fmt.Errorf("parse %q: %w", filePath, err)
			}

			instructions, err := parseOps(opsArgs)
			if err != nil {
				return err
			}

			result := patch.Apply(env, instructions)
			output := patch.Format(result)

			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
					return fmt.Errorf("write output: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Patched env written to %s\n", outputFile)
			} else {
				fmt.Fprint(cmd.OutOrStdout(), output)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write result to file instead of stdout")
	rootCmd.AddCommand(cmd)
}

func parseOps(args []string) ([]patch.Instruction, error) {
	var instructions []patch.Instruction
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid operation %q: expected op:arg", arg)
		}
		op, rest := patch.Op(strings.ToLower(parts[0])), parts[1]
		switch op {
		case patch.OpSet:
			kv := strings.SplitN(rest, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("set operation requires KEY=VALUE, got %q", rest)
			}
			instructions = append(instructions, patch.Instruction{Op: op, Key: kv[0], Value: kv[1]})
		case patch.OpDelete:
			instructions = append(instructions, patch.Instruction{Op: op, Key: rest})
		case patch.OpRename:
			kv := strings.SplitN(rest, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("rename operation requires OLD=NEW, got %q", rest)
			}
			instructions = append(instructions, patch.Instruction{Op: op, Key: kv[0], NewKey: kv[1]})
		default:
			return nil, fmt.Errorf("unknown operation %q: use set, delete, or rename", op)
		}
	}
	return instructions, nil
}
