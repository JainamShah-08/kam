package pipelines

import (
	"testing"

	"github.com/spf13/afero"
)

func TestBootstrapNew(t *testing.T) {
	type args struct {
		o     *GeneratorOptions
		appFs afero.Afero
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BootstrapNew(tt.args.o, tt.args.appFs); (err != nil) != tt.wantErr {
				t.Errorf("BootstrapNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
