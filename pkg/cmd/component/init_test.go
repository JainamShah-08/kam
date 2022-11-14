package bootstrapnew

import (
	"testing"

	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
)

var appFSTest = ioutils.NewMemoryFilesystem()

func Test_checkGit(t *testing.T) {

	dir1Full := "/app1/.git"
	dir1 := "/app1"
	dir2 := "/app2"

	appFSTest.Create(dir1Full)
	appFSTest.MkdirAll(dir2, 0755)
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			"Present",
			dir1,
			true,
		},
		{
			"Absent",
			dir2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkGit(appFSTest, tt.args); got != tt.want {
				t.Errorf("checkGit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkApplicationPath(t *testing.T) {
	dir1Full := "/app1/components"
	appFSTest.MkdirAll(dir1Full, 0755)
	dir2 := "/app2"
	dir1 := "/app1"

	// dir3Full := "/app3/apps/components"
	// dir3 := "/app3"
	// appFSTest.MkdirAll(dir3Full, 0755)

	tests := []struct {
		name    string
		args    string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			"Present",
			dir1,
			nil,
		},
		{
			"Absent",
			dir2,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckApplicationPath(appFSTest, tt.args); got != tt.wantErr {
				t.Errorf("CheckApplicationPath() = %v, want %v", got, tt.wantErr)
			}

		})
	}
}
