package addui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	PathGiven    string
	AppNameGiven string
)

func AddComponentName() string {
	var componentName string
	prompt := &survey.Input{
		Message: "Provide the Component name for your application ",
		Help:    "Required Field",
	}
	err := survey.AskOne(prompt, &componentName, validateCompNameAndPath())
	ui.HandleError(err)
	return strings.TrimSpace(componentName)
}
func validateCompNameAndPath() survey.Validator {
	return func(input interface{}) error {
		return validateCompNameAndPathFolder(input)
	}
}
func validateCompNameAndPathFolder(input interface{}) error {
	if u, ok := input.(string); ok {
		err := ui.ValidateName(u)
		if err != nil {
			return err
		}
		ui.HandleError(err)
		exists, anyErr := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(PathGiven, AppNameGiven, "components", u))
		if exists {
			return fmt.Errorf("the component : %s  already exists in path %s", u, filepath.Join(PathGiven, AppNameGiven, "components", u))
		}
		ui.HandleError(anyErr)
	}
	return nil
}
func AddApplicationNameComp() string {
	var applicationName string
	prompt := &survey.Input{
		Message: "Provide the Application name to add a Component",
		Help:    "Required Field",
	}
	err := survey.AskOne(prompt, &applicationName, validateNameAndPath())
	ui.HandleError(err)
	return strings.TrimSpace(applicationName)
}

func validateNameAndPath() survey.Validator {
	return func(input interface{}) error {
		return validateNameAndPathFolder(input)
	}
}
func validateNameAndPathFolder(input interface{}) error {
	if u, ok := input.(string); ok {
		err := ui.ValidateName(u)
		if err != nil {
			return err
		}
		ui.HandleError(err)
		exists, anyErr := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(PathGiven, u))
		if !exists {
			return fmt.Errorf("the given path : %s  doesnot exists ", filepath.Join(PathGiven, u))
		}
		ui.HandleError(anyErr)
	}
	return nil
}
func AddOutputPath() string {
	var path string
	prompt := &survey.Input{
		Message: "Provide a path to write where the Application is present",
		Help:    "This is the path where the GitOps repository configuration is stored locally before you push it to the repository GitRepoURL",
	}
	err := survey.AskOne(prompt, &path, validateOutputFolder())
	ui.HandleError(err)
	return strings.TrimSpace(path)
}

func validateOutputFolder() survey.Validator {
	return func(input interface{}) error {
		return validateOutput(input)
	}
}

func validateOutput(input interface{}) error {
	if u, ok := input.(string); ok {
		exists, err := ioutils.IsExisting(ioutils.NewFilesystem(), u)
		if !exists {
			return fmt.Errorf("the given path : %s  doesnot exists ", u)
		}
		ui.HandleError(err)
	}
	return nil
}
