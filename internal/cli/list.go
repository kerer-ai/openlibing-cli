package cli

import (
	"fmt"

	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/openlibing/openlibing-cli/pkg/spc"
	"github.com/spf13/cobra"
)

func newListCmd(reg *registry.RegistryImpl) *cobra.Command {
	var category string
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available Super Powers",
		Long:  "List all discoverable SPC (Super Power) definitions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var defs []*spc.SPCDefinition
			if category != "" {
				defs = reg.ListByCategory(category)
			} else {
				defs = reg.ListAll()
			}

			if len(defs) == 0 {
				fmt.Println("No Super Powers found.")
				return nil
			}

			switch format {
			case "json":
				rows := make([]map[string]interface{}, len(defs))
				for i, d := range defs {
					rows[i] = map[string]interface{}{
						"name":        d.Name,
						"type":        d.Type,
						"category":    d.Category,
						"description": d.Description,
						"origin":      d.Origin,
						"version":     d.Version,
					}
				}
				out, err := FormatJSON(rows)
				if err != nil {
					return err
				}
				fmt.Print(out)
			default:
				// Table format
				rows := make([]map[string]interface{}, len(defs))
				for i, d := range defs {
					rows[i] = map[string]interface{}{
						"name":        d.Name,
						"type":        d.Type,
						"category":    d.Category,
						"description": truncateString(d.Description, 60),
					}
				}
				fields := []spc.Field{
					{Name: "name", Header: "NAME"},
					{Name: "type", Header: "TYPE"},
					{Name: "category", Header: "CATEGORY"},
					{Name: "description", Header: "DESCRIPTION"},
				}
				fmt.Print(FormatTable(rows, fields))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category (pipeline|codecheck|sca)")
	cmd.Flags().StringVar(&format, "format", "table", "Output format (table|json)")

	return cmd
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
