package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"

	"github.com/openshift/odo/pkg/log"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"

	component "github.com/redhat-developer/kam/pkg/pipelines/component"
	hpipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	env "github.com/redhat-developer/kam/pkg/pipelines/component/environment"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	// AddEnvRecommendedCommandName the recommended command name
	AddEnvRecommendedCommandName = "add"
	componentNameFlag            = "component-name"
	outputFolderNameFlag         = "output"
	applicationNameFlag          = "application-name"
	environmentNameFlag          = "env-name"

	ApplicationNameAddEnvironmentMessage = "Provide the Application name to add an Environment"
)

var (
	addEnvExample = ktemplates.Examples(`
	# Add a new environment to the application's component in the GitOps repository
	# Example: kam env add --output <path to Application folder> --application-name <Application name> --component-name <component name> --env-name <environment name>
	
	%[1]s 
	`)

	addEnvLongDesc  = ktemplates.LongDesc(`Add a new environment to the application's component in the GitOps repository`)
	addEnvShortDesc = `Add a new environment`
)

// AddEnvParameters encapsulates the parameters for the kam pipelines init command.
type AddEnvParameters struct {
	*hpipelines.GeneratorOptions
	Interactive bool
}

// NewAddEnvParameters bootstraps a AddEnvParameters instance.
func NewAddEnvParameters() *AddEnvParameters {
	return &AddEnvParameters{
		GeneratorOptions: &hpipelines.GeneratorOptions{},
	}
}

// Checking the mandatory flags & reusing missingFlagErr in .go
func checkMandatoryFlags(flags map[string]string) error {
	missingFlags := []string{}
	mandatoryFlags := []string{componentNameFlag, applicationNameFlag, environmentNameFlag}
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

func initiateNonInteractiveModeComponent(io *AddEnvParameters) error {
	appFs := ioutils.NewFilesystem()
	mandatoryFlags := map[string]string{componentNameFlag: io.ComponentName, applicationNameFlag: io.ApplicationName, environmentNameFlag: io.EnvironmentName}
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
	if io.Output == "." {
		path, _ := os.Getwd()
		io.Output = path
	}
	exists, _ := ioutils.IsExisting(appFs, io.Output)
	if !exists {
		return fmt.Errorf("the provided path  %s does not exist. ", io.Output)
	}
	exists, _ = ioutils.IsExisting(appFs, filepath.Join(io.Output, io.ApplicationName))
	if !exists {
		return fmt.Errorf("the %s application doesn't exist in the path: %s", io.ApplicationName, io.Output)
	}
	exists, _ = ioutils.IsExisting(appFs, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
	if !exists {
		return fmt.Errorf("the %s component doesn't exist in path %s", io.ComponentName, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
	}
	exists, _ = ioutils.IsExisting(appFs, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName, "overlays", io.EnvironmentName))
	if exists {
		return fmt.Errorf("the %s environment already exists. Choose a different environment name", io.EnvironmentName)
	}
	return nil
}

func initiateInteractiveModeComponent(io *AddEnvParameters, cmd *cobra.Command) error {
	appFs := ioutils.NewFilesystem()
	log.Progressf("\nStarting interactive prompt\n")
	useDefaultValues := !ui.UseDefaultValuesComponent()
	if !cmd.Flag("output").Changed && useDefaultValues {
		io.Output = ui.ComponentOutputPath()
	}
	if io.Output == "./" {
		promp := !ui.UseDefaultValuesComponent()
		if promp {
			// ask for output folder
			io.Output = ui.ComponentOutputPath()
		}
	}
	if io.Output != "" {
		exists, _ := ioutils.IsExisting(appFs, io.Output)
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
			io.ApplicationName = ui.SelectApplicationNameComp(ApplicationNameAddEnvironmentMessage)
		} else {
			exists, _ := ioutils.IsExisting(appFs, filepath.Join(io.Output, io.ApplicationName))
			if !exists {
				log.Progressf("the Application : %s doesn't exists in Path %s", io.ApplicationName, io.Output)
				io.ApplicationName = ui.SelectApplicationNameComp(ApplicationNameAddEnvironmentMessage)
			}
		}
	} else {
		io.ApplicationName = ui.SelectApplicationNameComp(ApplicationNameAddEnvironmentMessage)
	}
	ui.AppNameGiven = io.ApplicationName
	presentComponents := ui.NumberOfComponents(filepath.Join(io.Output, io.ApplicationName, "components"))
	if len(presentComponents) == 0 {
		return fmt.Errorf("there are no components in the %s application in folder: %s ", io.ApplicationName, io.Output)
	} else {
		if io.ComponentName != "" {
			err := ui.ValidateName(io.ComponentName)
			if err != nil {
				log.Progressf("%v", err)
				io.ComponentName = ui.SelectComponentNameDelete(filepath.Join(io.Output, io.ApplicationName, "components"))
			}
			exists, _ := ioutils.IsExisting(appFs, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
			if !exists {
				log.Errorf("the component : %s does not exists in Application : %s at %s ", io.ComponentName, io.ApplicationName, io.Output)
				io.ComponentName = ui.SelectComponentNameDelete(filepath.Join(io.Output, io.ApplicationName, "components"))
			}
		} else {
			io.ComponentName = ui.SelectComponentNameDelete(filepath.Join(io.Output, io.ApplicationName, "components"))
		}
	}
	ui.ComponentNameGiven = io.ComponentName
	if io.EnvironmentName != "" {
		err := ui.ValidateName(io.EnvironmentName)
		if err != nil {
			log.Progressf("%v", err)
			io.EnvironmentName = ui.AddEnvironmentName()
		}
		exists, _ := ioutils.IsExisting(appFs, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName, "overlays", io.EnvironmentName))
		if exists {
			io.EnvironmentName = ui.AddEnvironmentName()
		}
	} else {
		io.EnvironmentName = ui.AddEnvironmentName()
	}

	return nil
}

// Complete completes AddEnvParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (eo *AddEnvParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 || eo.Interactive {
		return initiateInteractiveModeComponent(eo, cmd)
	}
	return initiateNonInteractiveModeComponent(eo)
}

// Validate validates the parameters of the EnvParameters.
func (eo *AddEnvParameters) Validate() error {
	return nil
}

// Run runs the project bootstrap command.
func (eo *AddEnvParameters) Run() error {
	appFs := ioutils.NewFilesystem()
	options := component.GeneratorOptions{
		ComponentName:   eo.ComponentName,
		ApplicationName: eo.ApplicationName,
		Output:          eo.Output,
		EnvironmentName: eo.EnvironmentName,
	}
	err := env.AddEnv(&options, appFs)
	if err != nil {
		return err
	}
	log.Successf("Created %s environment successfully. Check the deployment-patch.yaml to further customize.", eo.EnvironmentName)
	return nil
}

// NewCmdAddEnv creates the project add environment command.
func NewCmdAddEnv(name, fullName string) *cobra.Command {
	o := NewAddEnvParameters()

	addEnvCmd := &cobra.Command{
		Use:     name,
		Short:   addEnvShortDesc,
		Long:    addEnvLongDesc,
		Example: fmt.Sprintf(addEnvExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	addEnvCmd.Flags().StringVar(&o.EnvironmentName, "env-name", "", "Name of the environment/namespace")
	addEnvCmd.Flags().StringVar(&o.Output, "output", ".", "Folder path to the Application to add a new environment")
	addEnvCmd.Flags().StringVar(&o.ComponentName, "component-name", "", "Name of the component to add the environment")
	addEnvCmd.Flags().BoolVar(&o.Interactive, "interactive", false, "If true, enable prompting for most options if not already specified on the command line")
	addEnvCmd.Flags().StringVar(&o.ApplicationName, "application-name", "", "Name of the Application to add the environment")

	return addEnvCmd
}
