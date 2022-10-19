package environment

import (
	"fmt"

	"github.com/redhat-developer/gitops-generator/api/v1alpha1"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	component "github.com/redhat-developer/kam/pkg/pipelines/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/spf13/afero"
)

// AddEnv adds a new environment to the pipelines file.
func AddEnv(o *component.GeneratorOptions, appFs afero.Afero) error {
	bindingConfig := v1alpha1.GeneratorOptions{
		Name:     o.ComponentName,
		Replicas: 1,
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("256Mi"),
			},
		},
		OverlayEnvVar: nil,
	}

	e := gitops.NewCmdExecutor()
	anyErr := gitops.GenerateOverlaysAndPush(o.Output, false, "", bindingConfig, o.ApplicationName, o.EnvironmentName, "", "", e, appFs, "main", "", false, nil)
	if anyErr != nil {
		return fmt.Errorf("failed to create the environment :%s in component: %s: %w", o.EnvironmentName, o.ComponentName, anyErr)
	}
	return nil
}
