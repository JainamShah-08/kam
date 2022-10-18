package pipelines

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/redhat-developer/gitops-generator/api/v1alpha1"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/spf13/afero"
)

func AddComponent(o *pipelines.GeneratorOptions, appFs afero.Afero) error {
	componentSpec := v1alpha1.GeneratorOptions{
		Application: o.ApplicationName,
		Name:        o.ComponentName,
		TargetPort:  o.TargetPort,
		Route:       o.Route,
	}

	anyErr := gitops.GenerateAndPush(o.Output, "", componentSpec, gitops.NewCmdExecutor(), appFs, "main", false, "KAM cli")
	if anyErr != nil {
		return fmt.Errorf("failed to create the Component :%s in Application: %s: %w", o.ComponentName, o.ApplicationName, anyErr)
	}
	return nil
}

func DeleteComponent(o *pipelines.GeneratorOptions, appFs afero.Afero) error {
	anyErr := removeAndPush(filepath.Join(o.Output, o.ApplicationName), "", o.ComponentName, gitops.NewCmdExecutor(), appFs, "main", "", false, false)
	if anyErr != nil {
		return fmt.Errorf("failed to delete the Component :%s in Application: %s: %w", o.ComponentName, o.ApplicationName, anyErr)
	}
	return nil
}

func NewCmdExecutor() CmdExecutor {
	return CmdExecutor{}
}

type CmdExecutor struct {
}

func (e CmdExecutor) Execute(baseDir, command string, args ...string) ([]byte, error) {
	c := exec.Command(command, args...)
	c.Dir = baseDir
	output, err := c.CombinedOutput()
	return output, err
}

type Executor interface {
	Execute(baseDir, command string, args ...string) ([]byte, error)
}

func removeAndPush(outputPath string, remote string, componentName string, e Executor, appFs afero.Afero, branch string, context string, doPush bool, doClone bool) error {
	repoPath := filepath.Join(outputPath)
	if doClone {
		if out, err := e.Execute(outputPath, "git", "clone", remote, componentName); err != nil {
			return fmt.Errorf("failed to clone git repository in %q %q: %s", outputPath, string(out), err)
		}
		if _, err := e.Execute(repoPath, "git", "switch", branch); err != nil {
			if out, err := e.Execute(repoPath, "git", "checkout", "-b", branch); err != nil {
				return fmt.Errorf("failed to checkout branch %q in %q %q: %s", branch, repoPath, string(out), err)
			}
		}
	}
	// Generate the gitops resources and update the parent kustomize yaml file
	gitopsFolder := filepath.Join(repoPath, context)
	componentPath := filepath.Join(gitopsFolder, "components", componentName)
	if out, err := e.Execute(repoPath, "rm", "-rf", componentPath); err != nil {
		return fmt.Errorf("failed to delete %q folder in repository in %q %q: %s", componentPath, repoPath, string(out), err)
	}

	if doPush {
		return gitops.CommitAndPush(outputPath, "", remote, componentName, e, branch, fmt.Sprintf("Removed component %s", componentName))
	}

	return nil
}
