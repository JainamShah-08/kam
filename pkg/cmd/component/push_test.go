package bootstrapnew

import "testing"

func Test_checkapplicationPush(t *testing.T) {
	dir1Full := "/app1/components"
	dir1 := "/app1"
	dir2 := "/app2"

	appFSTest.Create(dir1Full)
	appFSTest.MkdirAll(dir2, 0755)
	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{"Correct Path", dir1, false},
		{"Incorrect Path", dir2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkapplicationPush(appFSTest, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("checkapplicationPush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
