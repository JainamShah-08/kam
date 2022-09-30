package pipelines

import (
	"path/filepath"
	"testing"

	hpipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
)

func TestAddAndDeleteComponent(t *testing.T) {
	fs := ioutils.NewFilesystem()
	testPath := "/user/example1"

	testOptions := []struct {
		output          string
		applicationName string
		componentName   string
		targetPort      int
		route           string
	}{
		{testPath, "app1", "comp1", 9090, "route1"},
		{testPath, "app2", "comp2", 8090, "route2"},
		{testPath, "app3", "comp2", 8090, "route3"},
	}
	for _, tt := range testOptions {
		path := fs.GetTempDir(filepath.Join(tt.output, tt.applicationName, "components"))
		x := len(path) - len(tt.applicationName) - len("components") - 3
		tt.output = path[0:x]
		o := &hpipelines.GeneratorOptions{
			Output:          tt.output,
			ApplicationName: tt.applicationName,
			ComponentName:   tt.componentName,
			TargetPort:      tt.targetPort,
			Route:           tt.route,
		}
		err := AddComponent(o, fs)
		if err != nil {
			t.Errorf("AddComponent() faced unexpected error %v", err)
		}

		exists, _ := fs.Exists(filepath.Join(tt.output, tt.applicationName, "components", tt.componentName))
		if !exists {
			t.Errorf("The path doesnot exist %s", tt.output)
		}
		err = DeleteComponent(o, fs)
		if err != nil {
			t.Errorf("DeleteComponent() faced unexpected error%v", err)
		}
		exists, _ = fs.Exists(filepath.Join(tt.output, tt.applicationName, "components", tt.componentName))
		if exists {
			t.Errorf("The path doesnot exist %s", tt.output)
		}

	}
}
