package component

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/log"

	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
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
	*pipelines.CompomemtParameters
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
	if io.OutputFolder == "./" {
		path, _ := os.Getwd()
		io.OutputFolder = path
	}
	exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), io.OutputFolder)
	if !exists {
		return fmt.Errorf("the given Path : %s  doesn't exists ", io.OutputFolder)
	}
	exists, _ = ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.OutputFolder, io.ApplicationName))
	if !exists {
		return fmt.Errorf("the given Application: %s  doesn't exists in the Path: %s", io.ApplicationName, io.OutputFolder)
	}
	exists, _ = ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.OutputFolder, io.ApplicationName, "components", io.ComponentName))
	if exists {
		return fmt.Errorf("the Component : %s  already exists in Path %s", io.ComponentName, filepath.Join(io.OutputFolder, io.ApplicationName, "components", io.ComponentName))
	}

	return nil
}

func initiateInteractiveModeComponent(io *AddCompParameters, cmd *cobra.Command) error {
	log.Progressf("\nStarting interactive prompt\n")
	promp := !ui.UseDefaultValuesComponent()
	if promp || cmd.Flag("output").Changed {
		// ask for output folder
		io.OutputFolder = ui.AddOutputPath()
	}
	if io.OutputFolder != "" {
		// Check for the path wether it is valid or not
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), io.OutputFolder)

		if !exists {
			log.Progressf("the provided Path doesn't exists in you directory : %s", io.OutputFolder)
			io.OutputFolder = ui.AddOutputPath()
			// ask for output folder
		}
	}
	if io.OutputFolder == "./" {
		path, _ := os.Getwd()
		io.OutputFolder = path
	}
	ui.PathGiven = io.OutputFolder
	if io.ApplicationName != "" {
		err := ui.ValidateName(io.ApplicationName)
		if err != nil {
			log.Progressf("%v", err)
			io.ApplicationName = ui.SelectApplicationNameComp()
		}
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.OutputFolder, io.ApplicationName))
		if !exists {
			log.Progressf("the Application : %s doesn't exists in Path %s", io.ApplicationName, io.OutputFolder)
			io.ApplicationName = ui.SelectApplicationNameComp()
		}
	}
	if io.ApplicationName == "" {
		io.ApplicationName = ui.SelectApplicationNameComp()
	}
	ui.AppNameGiven = io.ApplicationName

	if io.ComponentName != "" {
		err := ui.ValidateName(io.ComponentName)
		if err != nil {
			log.Progressf("%v", err)
			io.ComponentName = ui.SelectComponentNameComp()
		}
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.OutputFolder, io.ApplicationName, "components", io.ComponentName))
		if exists {
			log.Progressf("the component :%s already exists in Application : %s at %s ", io.ComponentName, io.ApplicationName, io.OutputFolder)
			io.ComponentName = ui.SelectComponentNameComp()
		}
	}
	if io.ComponentName == "" {
		io.ComponentName = ui.SelectComponentNameComp()
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

	err := pipelines.AddComponent(io.CompomemtParameters, appFs)
	if err != nil {
		return err
	}
	if err == nil {
		log.Successf("Created Component : %s in Application : %s at %s", io.ComponentName, io.ApplicationName, io.OutputFolder)
	}
	nextSteps()
	return nil
}

// NewAddCompParameters bootstraps a AddCompParameters instance.
func NewAddCompParameters() *AddCompParameters {
	return &AddCompParameters{
		CompomemtParameters: &pipelines.CompomemtParameters{},
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
	addCompCmd.Flags().StringVar(&o.OutputFolder, "output", "./", "Folder path to the Application to add the Component")
	addCompCmd.Flags().StringVar(&o.ApplicationName, "application-name", "", "Name of the Application to add a Component")
	return addCompCmd
}
func nextSteps() {
	log.Success("New Component added to the Application successfully\n\n",
		"\n",
	)
}
