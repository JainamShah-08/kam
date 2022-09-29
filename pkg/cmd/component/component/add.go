package component

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/log"

	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	hpipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	// AddCompRecommendedCommandName the recommended command name
	AddCompRecommendedCommandName = "add"
)

// AddCompParameters encapsulates the parameters for the kam pipelines init command.
type AddCompParameters struct {
	*hpipelines.GeneratorOptions
	Interactive bool
}

const (
	componentNameFlag    = "component-name"
	outputFolderNameFlag = "output"
	applicationNameFlag  = "application-name"
)

var (
	addCompExample = ktemplates.Examples(`
	# Add a new Component to Application
	# Example: kam component add --output <path to Application folder> --application-name <Application name to add component> --component-name new-component 
	
	%[1]s 
	`)

	addCompLongDesc  = ktemplates.LongDesc(`Add a new Component to the Application`)
	addCompShortDesc = `Add a new Component`
)

// Complete completes AddCompParameters after they've been created.
func (io *AddCompParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 || io.Interactive {
		return initiateInteractiveModeComponent(io, cmd)
	}
	return initiateNonInteractiveModeComponent(io)
}

// Checking the mandatory flags & reusing missingFlagErr in .go
func checkMandatoryFlags(flags map[string]string) error {
	missingFlags := []string{}
	mandatoryFlags := []string{componentNameFlag, applicationNameFlag}
	for _, flag := range mandatoryFlags {
		if flags[flag] == "" {
			missingFlags = append(missingFlags, fmt.Sprintf("%q", flag))
		}
	}
	if len(missingFlags) > 0 {
		return missingFlagErr(missingFlags)
	}
	return nil
}

func missingFlagErr(flags []string) error {
	return fmt.Errorf("required flag(s) %s not set", strings.Join(flags, ", "))
}

func initiateNonInteractiveModeComponent(io *AddCompParameters) error {
	mandatoryFlags := map[string]string{componentNameFlag: io.ComponentName, applicationNameFlag: io.ApplicationName}
	if err := checkMandatoryFlags(mandatoryFlags); err != nil {
		return err
	}
	err := ui.ValidateName(io.ApplicationName)
	if err != nil {
		return err
	}
	err = ui.ValidateName(io.ComponentName)
	if err != nil {
		return err
	}
	if io.Output == "./" {
		path, _ := os.Getwd()
		io.Output = path
	}
	exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), io.Output)
	if !exists {
		return fmt.Errorf("the given Path : %s  doesn't exists ", io.Output)
	}
	exists, _ = ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName))
	if !exists {
		return fmt.Errorf("the given Application: %s  doesn't exists in the Path: %s", io.ApplicationName, io.Output)
	}
	exists, _ = ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
	if exists {
		return fmt.Errorf("the Component : %s  already exists in Path %s", io.ComponentName, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
	}

	return nil
}

func initiateInteractiveModeComponent(io *AddCompParameters, cmd *cobra.Command) error {
	log.Progressf("\nStarting interactive prompt\n")
	promp := !ui.UseDefaultValuesComponent()
	if !cmd.Flag("output").Changed && promp {
		io.Output = ui.ComponentOutputPath()
	}
	if io.Output != "" {
		// Check for the path whether it is valid or not
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), io.Output)

		if !exists {
			log.Progressf("the provided Path doesn't exists in you directory : %s", io.Output)
			io.Output = ui.ComponentOutputPath()
			// ask for output folder
		}
	}
	if io.Output == "./" || io.Output == "." {
		path, _ := os.Getwd()
		io.Output = path
	}
	ui.PathGiven = io.Output
	if io.ApplicationName != "" {
		err := ui.ValidateName(io.ApplicationName)
		if err != nil {
			log.Progressf("%v", err)
			io.ApplicationName = ui.SelectApplicationNameComp("add")
		}
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName))
		if !exists {
			log.Progressf("the Application : %s doesn't exists in Path %s", io.ApplicationName, io.Output)
			io.ApplicationName = ui.SelectApplicationNameComp("add")
		}
	} else {
		io.ApplicationName = ui.SelectApplicationNameComp("add")
	}
	ui.AppNameGiven = io.ApplicationName

	if io.ComponentName != "" {
		err := ui.ValidateName(io.ComponentName)
		if err != nil {
			log.Progressf("%v", err)
			io.ComponentName = ui.SelectComponentNameComp()
		}
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
		if exists {
			log.Progressf("the component :%s already exists in Application : %s at %s ", io.ComponentName, io.ApplicationName, io.Output)
			io.ComponentName = ui.SelectComponentNameComp()
		}
	} else {
		io.ComponentName = ui.SelectComponentNameComp()
	}
	if !cmd.Flag("target-port").Changed && promp {
		io.TargetPort = ui.AddTargetPort()
	}
	if !cmd.Flag("route").Changed && promp {
		io.Route = ui.SelectRoute()
	}
	return nil
}

// Validate validates the parameters of the AddCompParameters.
func (io *AddCompParameters) Validate() error {
	return nil
}

// Run runs the project bootstrap command.
func (io *AddCompParameters) Run() error {
	log.Progressf("\nAdding the new component to the Application\n")
	appFs := ioutils.NewFilesystem()

	err := pipelines.AddComponent(io.GeneratorOptions, appFs)
	if err != nil {
		return err
	}

	if err == nil {
		log.Successf("Created Component : %s in Application : %s at %s", io.ComponentName, io.ApplicationName, io.Output)
	}
	nextSteps()
	return nil
}

// NewAddCompParameters bootstraps a AddCompParameters instance.
func NewAddCompParameters() *AddCompParameters {
	return &AddCompParameters{
		GeneratorOptions: &hpipelines.GeneratorOptions{},
	}
}
func NewCmdAddComp(name, fullName string) *cobra.Command {
	o := NewAddCompParameters()

	addCompCmd := &cobra.Command{
		Use:     name,
		Short:   addCompShortDesc,
		Long:    addCompLongDesc,
		Example: fmt.Sprintf(addCompExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	addCompCmd.Flags().StringVar(&o.ComponentName, "component-name", "", "Name of the component")
	addCompCmd.Flags().BoolVar(&o.Interactive, "interactive", false, "If true, enable prompting for most options if not already specified on the command line")
	addCompCmd.Flags().StringVar(&o.Output, "output", "./", "Folder path to the Application to add the Component")
	addCompCmd.Flags().StringVar(&o.ApplicationName, "application-name", "", "Name of the Application to add a Component")
	addCompCmd.Flags().IntVar(&o.TargetPort, "target-port", 8080, "Provide the Target Port for your Application")
	addCompCmd.Flags().StringVar(&o.Route, "route", "", "Provide the route to expose the component with. If provided, it will be referenced in the generated route.yaml")
	return addCompCmd
}
func nextSteps() {
	log.Success("New Component added to the Application successfully\n\n",
		"\n",
	)
}
