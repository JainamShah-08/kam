package bootstrapnew

import (
	"fmt"

	"github.com/openshift/odo/pkg/log"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	PushRecommendedCommandName = "push"
	applicationFolderFlagsP    = "application-folder"
	commitMessage              = "commit-message"
)

var (
	pushExample = ktemplates.Examples(`
    # Push command to push changes to git.
		kam push --application-folder <path to application> --commit-message <Message to commit>
		
    %[1]s 
    `)

	pushLongDesc  = ktemplates.LongDesc(`Pushing the application to Git repository.`)
	pushShortDesc = `Perform the Git add, commit and push commands for untracked changes.`
)

func NewPushParameters() *PushParameters {
	return &PushParameters{
		GeneratorOptions: &pipelines.GeneratorOptions{},
	}
}

type PushParameters struct {
	*pipelines.GeneratorOptions
	Interactive bool
}

func nonInteractiveModePush(io *PushParameters) error {
	mandatoryFlags := map[string]string{applicationFolderFlagsP: io.ApplicationFolder, commitMessage: io.CommitMessage}
	if err := checkMandatoryFlagsDescribe(mandatoryFlags); err != nil {
		return err
	}
	if err := checkApplicationPath(io.ApplicationFolder); err != nil {
		return err
	}
	if io.CommitMessage == "" {
		return fmt.Errorf("commit message is required to push repository to git")
	}
	return nil
}
func interactiveModePush(io *PushParameters) error {
	if io.ApplicationFolder == "" {
		io.ApplicationFolder = ui.ApplicationOutputPath()
	}
	if io.ApplicationFolder != "" {
		err := checkApplicationPath(io.ApplicationFolder)
		if err != nil {
			log.Progressf("%v", err)
			io.ApplicationFolder = ui.ApplicationOutputPath()
		}
	}
	if io.CommitMessage == "" {
		io.CommitMessage = ui.CommitMessage()
	}
	return nil
}
func (io *PushParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 || io.Interactive {
		return interactiveModePush(io)
	}

	return nonInteractiveModePush(io)
}
func (io *PushParameters) Validate() error {
	return nil
}
func (io *PushParameters) Run() error {

	isGit := checkGit(io.ApplicationFolder)

	if !isGit {
		return fmt.Errorf("no git repository has been initilaized to push")
	} else {
		e := gitops.NewCmdExecutor()
		if out, err := e.Execute(io.ApplicationFolder, "git", "add", "."); err != nil {
			return fmt.Errorf("failed to add components to repository in %q %q: %s", io.ApplicationFolder, string(out), err)
		}
		if out, err := e.Execute(io.ApplicationFolder, "git", "commit", "-m", io.CommitMessage); err != nil {
			return fmt.Errorf("failed to commit files to repository in %q %q: %s", io.ApplicationFolder, string(out), err)
		}
		if out, err := e.Execute(io.ApplicationFolder, "git", "push", "-u", "origin", "main"); err != nil {
			return fmt.Errorf("failed push remote to repository %q %q: %s", io.ApplicationFolder, string(out), err)
		}
	}
	return nil
}
func Push(name, fullName string) *cobra.Command {
	o := NewPushParameters()
	pushCmd := &cobra.Command{
		Use:     "push",
		Short:   pushShortDesc,
		Long:    pushLongDesc,
		Example: fmt.Sprintf(pushExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	pushCmd.Flags().StringVar(&o.ApplicationFolder, "application-folder", "", "Provode the path to the application")
	pushCmd.Flags().StringVar(&o.CommitMessage, "commit-message", "", "Provode a message to commit changes to repository")
	pushCmd.Flags().BoolVar(&o.Interactive, "interactive", false, "If true, enable prompting for most options if not already specified on the command line")
	return pushCmd
}
