package pipelines

import (
	"fmt"

	"github.com/redhat-developer/gitops-generator/api/v1alpha1"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/afero"
)

func AddComponent(o *pipelines.GeneratorOptions, appFs afero.Fs) error {
	componentSpec := v1alpha1.ComponentSpec{
		Application:   o.ApplicationName,
		ComponentName: o.ComponentName,

		Source: v1alpha1.ComponentSource{
			ComponentSourceUnion: v1alpha1.ComponentSourceUnion{
				GitSource: &v1alpha1.GitSource{
					URL: "",
				},
			},
		},
	}
	BootstrapNewVal := v1alpha1.Component{
		Spec: componentSpec,
	}
	e := gitops.NewCmdExecutor()
	anyErr := gitops.GenerateAndPush(o.Output, "", BootstrapNewVal, e, ioutils.NewFilesystem(), "main", false, "KAM cli", nil)
	if anyErr != nil {
		return fmt.Errorf("failed to create the Component :%s in Application: %s: %w", o.ComponentName, o.ApplicationName, anyErr)
	}
	return nil
}
