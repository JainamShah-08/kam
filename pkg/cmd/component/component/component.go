package component

import (
	"fmt"
	"os"

	"github.com/redhat-developer/kam/pkg/cmd/utility"
	"github.com/spf13/cobra"
)

//  CompRecommendedCommandName is the recommended Component command name.
const CompRecommendedCommandName = "component"

// NewCmdComp create a new component command
func NewCmdComp(name, fullName string) *cobra.Command {

	addCompCmd := NewCmdAddComp(AddCompRecommendedCommandName, utility.GetFullName(fullName, AddCompRecommendedCommandName))

	var compCmd = &cobra.Command{
		Use:   name,
		Short: "Manage component in application",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, AddCompRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	compCmd.Flags().AddFlagSet(addCompCmd.Flags())
	compCmd.AddCommand(addCompCmd)

	compCmd.Annotations = map[string]string{"command": "main"}
	return compCmd
}
