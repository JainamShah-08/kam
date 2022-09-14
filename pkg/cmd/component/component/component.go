package component

import (
	"fmt"
	"os"

	"github.com/redhat-developer/kam/pkg/cmd/utility"
	"github.com/spf13/cobra"
)

//  CompRecommendedCommandName is the recommended Component command name.
const CompRecommendedCommandName = "component"

var (
	deleteEx = "kam component delete --output <path to Application folder> --application-name <Application name to delete component> --component-name <name of the component to delete>"
	addEx    = "kam component add --output <path to Application folder> --application-name <Application name to add component> --component-name new-component"
)

// NewCmdComp create a new component command
func NewCmdComp(name, fullName string) *cobra.Command {

	addCompCmd := NewCmdAddComp(AddCompRecommendedCommandName, utility.GetFullName(fullName, AddCompRecommendedCommandName))
	deleteCompCmd := NewCmdDeleteComp(DeleteCompRecommendedCommandName, utility.GetFullName(fullName, DeleteCompRecommendedCommandName))
	var compCmd = &cobra.Command{
		Use:   name,
		Short: "Manage component in application",
		Example: fmt.Sprintf("%s\n%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, addEx, deleteEx),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
	compCmd.AddCommand(addCompCmd)
	compCmd.AddCommand(deleteCompCmd)
	compCmd.Annotations = map[string]string{"command": "main"}
	return compCmd
}
