package component

import (
	"fmt"
	"os"
	"path/filepath"

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
	DeleteCompRecommendedCommandName = "delete"
)

type DeleteCompParameters struct {
	*hpipelines.GeneratorOptions
	Interactive bool
}

var (
	deleteCompExample = ktemplates.Examples(`
	# Delete an existing Component from the Application
	# Example: kam component delete --output <path to Application folder> --application-name <Application name to delete component> --component-name <name of the component to delete>
	
	%[1]s 
	`)

	deleteCompLongDesc  = ktemplates.LongDesc(`Delete an existing Component from the Application`)
	deleteCompShortDesc = `Delete an existing Component`
)

func (io *DeleteCompParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 || io.Interactive {
		return initiateInteractiveModeDeleteComponent(io, cmd)
	}
	return initiateNonInteractiveModeDeleteComponent(io)
}

// checkMandatoryFlags is a function to check mandatory flags - defined in add.go
func initiateNonInteractiveModeDeleteComponent(io *DeleteCompParameters) error {
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
		return fmt.Errorf("the given Path : %s doesn't exists ", io.Output)
	}
	exists, _ = ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName))
	if !exists {
		return fmt.Errorf("the given Application: %s doesn't exists in the Path: %s", io.ApplicationName, io.Output)
	}
	presentComponents := ui.NumberOfComponents(filepath.Join(io.Output, io.ApplicationName, "components"))
	if len(presentComponents) == 0 {
		return fmt.Errorf("there are no components in the %s application in folder: %s ", io.ApplicationName, io.Output)
	}
	exists, _ = ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
	if !exists {
		return fmt.Errorf("the Component : %s does not exists in Path %s", io.ComponentName, filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
	}

	return nil
}

func initiateInteractiveModeDeleteComponent(io *DeleteCompParameters, cmd *cobra.Command) error {
	log.Progressf("\nStarting interactive prompt\n")
	if io.Output == "./" {
		promp := !ui.UseDefaultValuesComponent()
		if promp {
			// ask for output folder
			io.Output = ui.ComponentOutputPath()
		}
	}
	if io.Output != "" {
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), io.Output)
		if !exists {
			log.Progressf("the provided Path doesn't exists in you directory : %s", io.Output)
			io.Output = ui.ComponentOutputPath()
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
			io.ApplicationName = ui.SelectApplicationNameComp("delete")
		}
		exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName))
		if !exists {
			log.Progressf("the Application : %s doesn't exists in Path %s", io.ApplicationName, io.Output)
			io.ApplicationName = ui.SelectApplicationNameComp("delete")
		}
	} else {
		io.ApplicationName = ui.SelectApplicationNameComp("delete")
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
			exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.Output, io.ApplicationName, "components", io.ComponentName))
			if !exists {
				log.Errorf("the component : %s does not exists in Application : %s at %s ", io.ComponentName, io.ApplicationName, io.Output)
				io.ComponentName = ui.SelectComponentNameDelete(filepath.Join(io.Output, io.ApplicationName, "components"))
			}
		} else {
			io.ComponentName = ui.SelectComponentNameDelete(filepath.Join(io.Output, io.ApplicationName, "components"))
		}
	}

	return nil
}

// Validate validates the parameters of the DeleteCompParameters.
func (io *DeleteCompParameters) Validate() error {
	return nil
}

// Run runs the project component delete command command.
func (io *DeleteCompParameters) Run() error {
	log.Progressf("\nDeleted the component from the Application\n")
	appFs := ioutils.NewFilesystem()

	err := pipelines.DeleteComponent(io.GeneratorOptions, appFs)
	if err != nil {
		return err
	}
	if err == nil {
		log.Successf("Deleted Component : %s in Application : %s at %s", io.ComponentName, io.ApplicationName, io.Output)
	}
	nextStepsDelete()
	return nil
}
func NewDeleteCompParameters() *DeleteCompParameters {
	return &DeleteCompParameters{
		GeneratorOptions: &hpipelines.GeneratorOptions{},
	}
}
func NewCmdDeleteComp(name, fullName string) *cobra.Command {
	o := NewDeleteCompParameters()

	deleteCompCmd := &cobra.Command{
		Use:     name,
		Short:   deleteCompShortDesc,
		Long:    deleteCompLongDesc,
		Example: fmt.Sprintf(deleteCompExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	deleteCompCmd.Flags().StringVar(&o.ComponentName, "component-name", "", "Name of the component to delete")
	deleteCompCmd.Flags().BoolVar(&o.Interactive, "interactive", false, "If true, enable prompting for most options if not already specified on the command line")
	deleteCompCmd.Flags().StringVar(&o.Output, "output", "./", "Folder path to the Application to delete the Component")
	deleteCompCmd.Flags().StringVar(&o.ApplicationName, "application-name", "", "Name of the Application to delete a Component")
	return deleteCompCmd
}
func nextStepsDelete() {
	log.Success("Component deleted from the Application successfully\n\n",
		"\n",
	)
}
