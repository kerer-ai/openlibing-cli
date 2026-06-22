package cli

import (
	"fmt"

	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

func newChatCmd(eng *engine.Engine, reg *registry.RegistryImpl) *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "Enter AI conversational mode (coming soon)",
		Long: `chat launches an interactive AI-powered REPL for querying OpenLibing data.

In chat mode, you use natural language to describe what you want to query.
The AI automatically selects the right Super Power, fills in parameters,
and explains the results.

Note: Chat mode requires LLM API credentials in ~/.openlibing/auth.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Chat mode is coming in a future release.")
			fmt.Println("For now, use 'openlibing run <spc-name>' for direct queries.")
			fmt.Println("Run 'openlibing list' to see available Super Powers.")
			return nil
		},
	}
}
