package environment

import (
	component "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"testing"
)

func TestGenerateDeploymentPatch(t *testing.T) {
	var appFs = ioutils.NewMemoryFilesystem()
	var testPath = "/fake-path/test"

	replicas := int32(1)
	image := "container-image"

	tests := []struct {
		name             string
		generatorOptions component.GeneratorOptions
		wantDeployment   appsv1.Deployment
		wantExist        bool
	}{
		{
			name: "Simple component, no optional fields set",
			generatorOptions: component.GeneratorOptions{
				ComponentName:   "frontend",
				ApplicationName: "testapp",
				Output:          testPath,
				EnvironmentName: "env-name",
			},
			wantExist: true,
			wantDeployment: appsv1.Deployment{
				TypeMeta: v1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &v1.LabelSelector{},
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "container-image",
									Image: image,
									Resources: corev1.ResourceRequirements{
										Limits: corev1.ResourceList{
											corev1.ResourceCPU:            resource.MustParse("1"),
											corev1.ResourceRequestsMemory: resource.MustParse("256Mi"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddEnv(&tt.generatorOptions, appFs)
			if err != nil {
				t.Errorf("Error adding environment. Got error: %v ", err)
			}

			pathOfEnvFolder := filepath.Join(testPath, tt.generatorOptions.ApplicationName, "components", tt.generatorOptions.ComponentName, "overlays", tt.generatorOptions.EnvironmentName)
			envExists, _ := appFs.Exists(pathOfEnvFolder)
			if envExists != tt.wantExist {
				t.Errorf("Error expect file to exist: %s ", pathOfEnvFolder)
			}

			pathOfDeploymentPatchFile := filepath.Join(testPath, tt.generatorOptions.ApplicationName, "components", tt.generatorOptions.ComponentName, "overlays", tt.generatorOptions.EnvironmentName, "deployment-patch.yaml")
			fileExists, _ := appFs.Exists(pathOfDeploymentPatchFile)
			if fileExists != tt.wantExist {
				t.Errorf("Error expect file to exist: %s ", pathOfDeploymentPatchFile)

			}
		})
	}
}
