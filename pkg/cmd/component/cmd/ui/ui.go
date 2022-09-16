package ui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/afero"
	"gopkg.in/AlecAivazis/survey.v1"
)

// EnterGitRepo allows the user to specify the git repository in a prompt.
func EnterGitRepoURL() string {
	var gitRepoURL string
	prompt := &survey.Input{
		Message: "Provide the URL for your GitOps repository",
		Help:    "The GitOps repository stores your GitOps configuration files, including your Openshift Pipelines resources for driving automated deployments and builds.  Please enter a valid git repository e.g. https://github.com/example/myorg.git",
	}
	err := survey.AskOne(prompt, &gitRepoURL, makeURLValidatorCheck())
	handleError(err)
	return strings.TrimSpace(gitRepoURL)
}

func VerifyOutput(appFs afero.Fs, originalPath string, overwrite bool, appName string, outputPathOverridden bool, promptForPath bool) (string, bool) {
	var outputPath = originalPath
	var doOverwrite = overwrite
	prompt := &survey.Input{
		Message: "Provide a path to write GitOps resources?",
		Help:    "This is the path where the GitOps repository configuration is stored locally before you push it to the repository GitRepoURL",
		Default: originalPath,
	}
	if !outputPathOverridden && promptForPath {
		handleError(survey.AskOne(prompt, &outputPath, nil))
		outputPath = strings.TrimSpace(outputPath)
	}
	for {
		exists, err := ioutils.IsExisting(appFs, filepath.Join(outputPath, appName))
		handleError(err)
		if !exists || overwrite {
			break
		}
		doOverwrite = SelectOptionOverwrite(outputPath)
		if doOverwrite {
			break
		}
		handleError(survey.AskOne(prompt, &outputPath, nil))
		outputPath = strings.TrimSpace(outputPath)
	}
	return outputPath, doOverwrite
}

func PathExists(appFs afero.Fs, path string) bool {
	exists, err := ioutils.IsExisting(appFs, path)
	handleError(err)
	return exists
}

// Not validating the token
func EnterGitSecret(repoURL string) string {
	var gitWebhookSecret string
	prompt := &survey.Password{
		Message: fmt.Sprintf("Provide a token used to authenticate requests to %s", repoURL),
		Help:    "Tokens are required to authenticate to git provider various operations on git repository (e.g. enable automated creation/push to git-repo).",
	}

	err := survey.AskOne(prompt, &gitWebhookSecret, makeSecretValidator())
	handleError(err)
	return gitWebhookSecret
}

// SelectOptionOverwrite allows users the option to overwrite the current gitops configuration locally through the UI prompt.
func SelectOptionOverwrite(currentPath string) bool {
	var overwrite string
	prompt := &survey.Select{
		Message: "Do you want to overwrite your output path?",
		Help:    "Overwrite: " + currentPath,
		Options: []string{"yes", "no"},
		Default: "no",
	}
	handleError(survey.AskOne(prompt, &overwrite, nil))
	return overwrite == "yes"
}

// SelectPrivateRepoDriver lets users choose the driver for their git hosting
// service.
func SelectPrivateRepoDriver() string {
	var driver string
	prompt := &survey.Select{
		Message: "Please select which driver to use for your Git host",
		Options: []string{"github", "gitlab"},
	}

	err := survey.AskOne(prompt, &driver, survey.Required)
	handleError(err)
	return driver
}

// SelectOptionPushToGit allows users the option to select if they
// want to incorporate the feature of the commit status tracker through the UI prompt.
func SelectOptionPushToGit() bool {
	var optionPushToGit string
	prompt := &survey.Select{
		Message: "Do you want to create and push the resources to your gitops repository?",
		Help:    "This will create a private repository, commit and push the generated resources and requires an auth token with the correct privileges",
		Options: []string{"yes", "no"},
	}
	err := survey.AskOne(prompt, &optionPushToGit, survey.Required)
	handleError(err)
	return optionPushToGit == "yes"
}

// UseDefaultValues allows users to use default values so that they will be prompted with fewer questions in interactive mode
func UseDefaultValues() bool {
	var defaultFlagVal = map[string]string{
		"output":    "\"./\"",
		"namespace": "openshift-gitops",
	}
	flagValues := "\n\nThe default values used for the options, if not overwritten from the command line, are:\n"
	buff := bytes.Buffer{}
	w := tabwriter.NewWriter(&buff, 0, 8, 1, '\t', tabwriter.AlignRight)
	for k, v := range defaultFlagVal {
		fmt.Fprintf(w, "--%v\t%v\n", k, v)
	}
	w.Flush()
	vStr := buff.String()
	flagValues += vStr

	var useDefaults string
	prompt := &survey.Select{
		Message: "Do you want to accept all default values and be prompted only for the minimum required options?",
		Help:    "Select yes to accept default values or select no to be prompted for all options that haven't already been specified on the command line" + flagValues,
		Options: []string{"yes", "no"},
		Default: "yes",
	}
	handleError(survey.AskOne(prompt, &useDefaults, nil))
	return useDefaults == "yes"
}

func AddApplicationName() string {
	var applicationName string
	prompt := &survey.Input{
		Message: "Provide the Application name ",
		Help:    "Required Field",
	}
	err := survey.AskOne(prompt, &applicationName, MakeNameCheck())
	handleError(err)
	return strings.TrimSpace(applicationName)
}

func AddComponentName() string {
	var componentName string
	prompt := &survey.Input{
		Message: "Provide the Component name for your application ",
		Help:    "Required Field",
	}
	err := survey.AskOne(prompt, &componentName, MakeNameCheck())
	handleError(err)
	return strings.TrimSpace(componentName)
}

func MakeNameCheck() survey.Validator {
	return func(input interface{}) error {
		return validateName(input)
	}
}

func validateName(input interface{}) error {
	if s, ok := input.(string); ok {
		err := ValidateName(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddTargetPort() int {
	var targetPort int
	prompt := &survey.Input{
		Message: "Provide the Target Port ",
	}
	err := survey.AskOne(prompt, &targetPort, makeTargetPortCheck())
	handleError(err)
	return targetPort
}

func makeTargetPortCheck() survey.Validator {
	return func(input interface{}) error {
		return validateTarget(input)
	}
}

func validateTarget(input interface{}) error {
	if s, ok := input.(int); ok {
		err := ValidateTargetPort(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// UseKeyringRingSvc , allows users an option between the Internal image registry and the external image registry through the UI prompt.
func UseKeyringRingSvc() bool {
	var optionImageRegistry string
	prompt := &survey.Select{
		Message: "Do you wish to securely store the git-host-access-token in the keyring on your local machine? (The token is saved and not validated.)",
		Help:    "The token will be stored securely in the keyring of your local mahine. It will be reused by kam commands(bootstrap/webhoook), further iteration of these commands will not prompt for the access-token",
		Options: []string{"yes", "no"},
	}

	err := survey.AskOne(prompt, &optionImageRegistry, survey.Required)
	handleError(err)
	return optionImageRegistry == "yes"
}
func SelectRoute() string {
	var optionRoute string
	prompt := &survey.Input{
		Message: "Provide a route name for your application ?",
		Help:    "If you specify the route flag and pass the string, that string will be in the route.yaml that is generated",
	}

	err := survey.AskOne(prompt, &optionRoute, nil)
	handleError(err)
	return strings.TrimSpace(optionRoute)
}

// ----------------------------------------- UI For Component ADD/Delete Command -----------------------------------------

var (
	PathGiven    string
	AppNameGiven string
)

func UseDefaultValuesComponent() bool {
	var useDefaults string
	prompt := &survey.Select{
		Message: "Do you want to accept all default values and be prompted only for the minimum required options?",
		Help:    "Select yes to accept default values or select no to be prompted for all options that haven't already been specified on the command line",
		Options: []string{"yes", "no"},
		Default: "yes",
	}
	handleError(survey.AskOne(prompt, &useDefaults, nil))
	return useDefaults == "yes"
}
func SelectComponentNameComp() string {
	var componentName string
	prompt := &survey.Input{
		Message: "Provide the Component name for your Application ",
		Help:    "Required Field",
	}
	err := survey.AskOne(prompt, &componentName, validateCompNameAndPath())
	handleError(err)
	return strings.TrimSpace(componentName)
}

// temperoray
func SelectComponentNameDelete(path string) string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var directory []string
	for _, f := range files {
		err = validateName(f.Name())
		if err == nil {
			directory = append(directory, f.Name())
		}
	}
	var componentName string
	prompt := &survey.Select{
		Message: "Provide the Component name for your Application ",
		Help:    "Required Field",
		Options: directory,
	}
	// err = survey.AskOne(prompt, &componentName, validateCompNameAndPath())
	// handleError(err)
	handleError(survey.AskOne(prompt, &componentName, nil))
	return strings.TrimSpace(componentName)
}

//exit
func validateCompNameAndPath() survey.Validator {
	return func(input interface{}) error {
		return validateCompNameAndPathFolder(input)
	}
}
func validateCompNameAndPathFolder(input interface{}) error {
	if u, ok := input.(string); ok {
		err := ValidateName(u)
		if err != nil {
			return err
		}
		handleError(err)
		exists, anyErr := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(PathGiven, AppNameGiven, "components", u))
		if exists {
			return fmt.Errorf("the Component : %s already exists in Application : %s at Path %s", u, AppNameGiven, PathGiven)
		}
		handleError(anyErr)
	}
	return nil
}
func SelectApplicationNameComp() string {
	var applicationName string
	prompt := &survey.Input{
		Message: "Provide the Application name to add/delete a Component",
		Help:    "Required Field",
	}
	err := survey.AskOne(prompt, &applicationName, validateNameAndPath())
	handleError(err)
	return strings.TrimSpace(applicationName)
}

func validateNameAndPath() survey.Validator {
	return func(input interface{}) error {
		return validateNameAndPathFolder(input)
	}
}
func validateNameAndPathFolder(input interface{}) error {
	if u, ok := input.(string); ok {
		err := ValidateName(u)
		if err != nil {
			return err
		}
		handleError(err)
		exists, anyErr := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(PathGiven, u))
		if !exists {
			return fmt.Errorf("the given Application : %s doesn't exists in the Path : %s", u, PathGiven)
		}
		handleError(anyErr)
	}
	return nil
}
func ComponentOutputPath() string {
	var path string
	prompt := &survey.Input{
		Message: "Provide a Path to write where the Application is present",
		Help:    "This is the path where the GitOps repository configuration is stored locally before you push it to the repository GitRepoURL",
	}
	err := survey.AskOne(prompt, &path, validateOutputFolder())
	handleError(err)
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
			return fmt.Errorf("the given Path : %s doesn't exists in your directory", u)
		}
		handleError(err)
	}
	return nil
}
