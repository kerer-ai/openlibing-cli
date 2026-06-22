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
		Long: `Execute a Super Power by name. Parameters are passed as flags (--param-name value).

For custom SPCs with unknown flags, pass them after a -- separator:
  openlibing run my-spc --project-id 123 -- --custom-flag value`,
		Args: cobra.ExactArgs(1),
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spcName := args[0]

			// Collect all known flags as parameters
			params := make(map[string]interface{})
			flags := cmd.Flags()
			flags.Visit(func(f *pflag.Flag) {
				if f.Name == "output" || f.Name == "param" {
					return
				}
				paramName := strings.ReplaceAll(f.Name, "-", "_")
				params[paramName] = f.Value.String()
			})

			// Collect --param key=value pairs (for custom SPC parameters)
			paramFlags, _ := cmd.Flags().GetStringArray("param")
			for _, p := range paramFlags {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) == 2 {
					params[parts[0]] = parts[1]
				}
			}

			// Also collect unknown args from after --
			if cmd.Flags().ArgsLenAtDash() >= 0 {
				dashArgs := cmd.Flags().Args()[cmd.Flags().ArgsLenAtDash():]
				for i := 0; i < len(dashArgs)-1; i += 2 {
					name := strings.TrimPrefix(dashArgs[i], "--")
					name = strings.ReplaceAll(name, "-", "_")
					params[name] = dashArgs[i+1]
				}
			}

			result, err := eng.Execute(spcName, params)
			if err != nil {
				return err
			}

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

	// Pre-register common parameter flags
	cmd.Flags().String("project-id", "", "Project ID")
	cmd.Flags().String("run-id", "", "Pipeline run ID")
	cmd.Flags().String("pipeline-run-id", "", "Pipeline run ID")
	cmd.Flags().String("job-run-id", "", "Job run ID")
	cmd.Flags().String("step-run-id", "", "Step run ID")
	cmd.Flags().String("status", "", "Filter by status")
	cmd.Flags().Int("limit", 0, "Max results")
	// --param key=value for custom SPC parameters
	cmd.Flags().StringArrayP("param", "p", nil, "Pass custom parameter as key=value (can be repeated)")

	return cmd
}
