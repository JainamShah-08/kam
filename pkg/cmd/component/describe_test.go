package bootstrapnew

import (
	"testing"
)

// Test case for checking the input flags.
func TestMissingFlagsDescribe(t *testing.T) {
	tests := []struct {
		desc  string
		flags map[string]string
		err   error
	}{
		{
			"Required flags are present",
			map[string]string{"application-name": "value-1"},
			nil,
		},

		{
			"A required flag is absent",
			map[string]string{"application-name": ""},
			missingFlagErr([]string{`"application-name"`}),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			gotErr := checkMandatoryFlagsDescribe(test.flags)
			if gotErr != nil && test.err != nil {
				if gotErr.Error() != test.err.Error() {
					t.Fatalf("error mismatch: got %v, want %v", gotErr, test.err)
				}
			} else if gotErr != test.err {
				t.Fatalf("error mismatch: got %v, want %v", gotErr, test.err)
			}
		})
	}
}
