package pipelines

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/log"
	"github.com/redhat-developer/gitops-generator/api/v1alpha1"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/spf13/afero"
)

// function to populate the GitSource struct
func BootstrapNew(o *GeneratorOptions, appFs afero.Afero) error {

	genOptions := v1alpha1.GeneratorOptions{
		Application: o.ApplicationName,
		Name:        o.ComponentName,
		Secret:      o.Token,
		TargetPort:  o.TargetPort,
		Route:       o.Route,
		GitSource: &v1alpha1.GitSource{
			URL: o.GitRepoURL,
		},
	}

	if ui.PathExists(appFs, filepath.Join(o.Output, o.ApplicationName)) && !o.Overwrite {
		return fmt.Errorf("%v the application name already exists in given directory %v", o.ApplicationName, o.Output)
	}
	anyErr := gitops.NewGitopsGen().GenerateAndPush(o.Output, o.GitRepoURL, genOptions, appFs, "main", o.PushToGit, "KAM cli")
	if anyErr != nil {
		log.Progressf(o.GitRepoURL)
		return fmt.Errorf("failed to create the gitops repository: %q: %w", o.GitRepoURL, anyErr)
	}
	return nil
}
