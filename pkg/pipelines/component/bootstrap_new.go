package pipelines

import (
	"fmt"
	"path/filepath"

	"github.com/redhat-developer/gitops-generator/api/v1alpha1"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/spf13/afero"
)

// function to populate the GitSource struct
func BootstrapNew(o *GeneratorOptions, appFs afero.Afero) error {

	componentSpec := v1alpha1.ComponentSpec{
		Application:   o.ApplicationName,
		ComponentName: o.ComponentName,
		Secret:        o.Secret,
		TargetPort:    o.TargetPort,
		Route:         o.Route,

		Source: v1alpha1.ComponentSource{
			ComponentSourceUnion: v1alpha1.ComponentSourceUnion{
				GitSource: &v1alpha1.GitSource{
					URL: o.GitRepoURL,
				},
			},
		},
	}
	BootstrapNewVal := v1alpha1.Component{
		Spec: componentSpec,
	}

	if ui.PathExists(appFs, filepath.Join(o.Output, o.ApplicationName)) && !o.Overwrite {
		return fmt.Errorf("%v the application name already exists in given directory %v", o.ApplicationName, o.Output)
	}
	e := gitops.NewCmdExecutor()
	anyErr := gitops.GenerateAndPush(o.Output, o.GitRepoURL, BootstrapNewVal, e, appFs, "main", o.PushToGit, "KAM cli", nil)
	if anyErr != nil {
		return fmt.Errorf("failed to create the gitops repository: %q: %w", o.GitRepoURL, anyErr)
	}
	return nil
}
