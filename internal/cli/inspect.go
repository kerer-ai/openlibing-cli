package cli

import (
	"fmt"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

func newInspectCmd(reg *registry.RegistryImpl) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <spc-name>",
		Short: "Show full definition of a Super Power",
		Long:  "Display the complete SPC definition including parameters, output fields, examples, and AI hints.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def, err := reg.Get(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Name:        %s\n", def.Name)
			fmt.Printf("Version:     %s\n", def.Version)
			fmt.Printf("Type:        %s\n", def.Type)
			fmt.Printf("Category:    %s\n", def.Category)
			fmt.Printf("Origin:      %s\n", def.Origin)
			fmt.Printf("Tags:        %s\n", strings.Join(def.Tags, ", "))
			fmt.Printf("\nDescription:\n  %s\n", def.Description)

			if len(def.Parameters) > 0 {
				fmt.Printf("\nParameters:\n")
				for _, p := range def.Parameters {
					required := ""
					if p.Required {
						required = " (required)"
					}
					defaultVal := ""
					if p.Default != nil {
						defaultVal = fmt.Sprintf(" [default: %v]", p.Default)
					}
					fmt.Printf("  --%s  %s%s%s\n", toFlagName(p.Name), p.Type, required, defaultVal)
					if p.Description != "" {
						fmt.Printf("        %s\n", p.Description)
					}
					if len(p.Enum) > 0 {
						fmt.Printf("        Values: %s\n", strings.Join(p.Enum, ", "))
					}
				}
			}

			fmt.Printf("\nSource:\n")
			fmt.Printf("  %s %s\n", def.Source.Method, def.Source.Endpoint)

			fmt.Printf("\nOutput: %s", def.Output.Format)
			if len(def.Output.Fields) > 0 {
				fmt.Printf(" (%d fields)\n", len(def.Output.Fields))
				for _, f := range def.Output.Fields {
					fmt.Printf("  %-15s → %s", f.Name, f.Header)
					if f.Transform != "" {
						fmt.Printf(" [%s]", f.Transform)
					}
					fmt.Println()
				}
			} else {
				fmt.Println()
			}

			if len(def.Examples) > 0 {
				fmt.Printf("\nExamples:\n")
				for _, ex := range def.Examples {
					fmt.Printf("  %s\n", ex.Command)
					fmt.Printf("  # %s\n\n", ex.Description)
				}
			}

			return nil
		},
	}
}

// toFlagName converts parameter name to CLI flag format.
func toFlagName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}
