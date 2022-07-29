package component

import "testing"

// Test case for checking the input flags.
func TestMissingFlagsComponent(t *testing.T) {
	tests := []struct {
		desc  string
		flags map[string]string
		err   error
	}{
		{
			"Required flags are present",
			map[string]string{"component-name": "value-1", "application-name": "value-2", "output": "value-3"},
			nil,
		},
		{
			"A required flag is absent",
			map[string]string{"component-name": "", "application-name": "value-2", "output": "value-3"},
			missingFlagErr([]string{`"component-name"`}),
		},
		{
			"A required flag is absent",
			map[string]string{"component-name": "value-1", "application-name": "", "output": "value-3"},
			missingFlagErr([]string{`"application-name"`}),
		},
		{
			"A required flag is absent",
			map[string]string{"component-name": "value-1", "application-name": "value-2", "output": ""},
			missingFlagErr([]string{`"output"`}),
		},
		{
			"Multiple required flags are absent",
			map[string]string{"component-name": "", "application-name": "", "output": ""},
			missingFlagErr([]string{`"component-name"`, `"output"`, `"application-name"`}),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			gotErr := checkMandatoryFlags(test.flags)
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
