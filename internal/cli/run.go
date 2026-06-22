package cli

import (
	"fmt"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newRunCmd(eng *engine.Engine) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <spc-name>",
		Short: "Execute a Super Power",
		Long:  "Execute a Super Power by name. Parameters are passed as flags (--param-name value).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spcName := args[0]

			// Collect all non-standard flags as parameters
			params := make(map[string]interface{})
			flags := cmd.Flags()
			flags.Visit(func(f *pflag.Flag) {
				if f.Name == "output" {
					return // skip global flags
				}
				paramName := strings.ReplaceAll(f.Name, "-", "_")
				params[paramName] = f.Value.String()
			})

			result, err := eng.Execute(spcName, params)
			if err != nil {
				return err
			}

			// Override format if --output flag was set
			if outputFormat != "" {
				result.Format = outputFormat
			}

			output, err := FormatResult(result)
			if err != nil {
				return err
			}

			fmt.Print(output)
			return nil
		},
	}

	// Pre-register common parameter flags.
	// The engine validates parameters; we accept arbitrary flags here.
	cmd.Flags().String("project-id", "", "Project ID")
	cmd.Flags().String("run-id", "", "Pipeline run ID")
	cmd.Flags().String("pipeline-run-id", "", "Pipeline run ID")
	cmd.Flags().String("job-run-id", "", "Job run ID")
	cmd.Flags().String("step-run-id", "", "Step run ID")
	cmd.Flags().String("status", "", "Filter by status")
	cmd.Flags().Int("limit", 0, "Max results")

	return cmd
}
