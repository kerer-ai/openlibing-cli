package cli

import (
	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

// NewRootCmd creates the root command with all subcommands.
func NewRootCmd(eng *engine.Engine, reg *registry.RegistryImpl) *cobra.Command {
	root := &cobra.Command{
		Use:   "openlibing",
		Short: "OpenLibing CLI — AI-native CI/CD data query tool",
		Long: `openlibing is a CLI tool for querying OpenLibing CI/CD platform data.

All capabilities are defined as SPC (Skill Pipeline Configuration) files.
Use 'openlibing list' to see available Super Powers.
Use 'openlibing run <name>' to execute a Super Power.
Use 'openlibing inspect <name>' to view SPC details.
Use 'openlibing chat' to enter AI conversational mode.`,
		SilenceUsage: true,
	}

	root.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (table|json|yaml|raw)")

	root.AddCommand(newRunCmd(eng))
	root.AddCommand(newListCmd(reg))
	root.AddCommand(newInspectCmd(reg))
	root.AddCommand(newChatCmd(eng, reg))

	return root
}
