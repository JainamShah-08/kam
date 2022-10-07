package env

import (
	"fmt"
	"os"

	"github.com/redhat-developer/kam/pkg/cmd/utility"
	"github.com/spf13/cobra"
)

// EnvRecommendedCommandName is the recommended environment command name.
const EnvRecommendedCommandName = "env"

// NewCmdEnv create a new environment command
func NewCmdEnv(name, fullName string) *cobra.Command {

	addEnvCmd := NewCmdAddEnv(AddEnvRecommendedCommandName, utility.GetFullName(fullName, AddEnvRecommendedCommandName))
	var envCmd = &cobra.Command{
		Use:   name,
		Short: "Manage an environment in GitOps",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, "kam env add --output <path to Application folder> --application-name <Application name> --component-name <component name> --env-name <environment name>"),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
	envCmd.AddCommand(addEnvCmd)
	envCmd.Annotations = map[string]string{"command": "main"}
	return envCmd
}
