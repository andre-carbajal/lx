package parser

import (
	"testing"
)

func TestExtractGlobalFlags(t *testing.T) {
	tests := []struct {
		name              string
		input             []string
		expectedDry       bool
		expectedVerbose   bool
		expectedRemaining []string
	}{
		{
			name:              "no global flags",
			input:             []string{"ls", "-la", "/home"},
			expectedDry:       false,
			expectedVerbose:   false,
			expectedRemaining: []string{"ls", "-la", "/home"},
		},
		{
			name:              "only --dry",
			input:             []string{"--dry", "ls", "-la"},
			expectedDry:       true,
			expectedVerbose:   false,
			expectedRemaining: []string{"ls", "-la"},
		},
		{
			name:              "only --verbose",
			input:             []string{"--verbose", "ls", "-la"},
			expectedDry:       false,
			expectedVerbose:   true,
			expectedRemaining: []string{"ls", "-la"},
		},
		{
			name:              "both flags",
			input:             []string{"--dry", "--verbose", "ls", "-la"},
			expectedDry:       true,
			expectedVerbose:   true,
			expectedRemaining: []string{"ls", "-la"},
		},
		{
			name:              "-v shorthand for verbose",
			input:             []string{"-v", "ls", "-la"},
			expectedDry:       false,
			expectedVerbose:   true,
			expectedRemaining: []string{"ls", "-la"},
		},
		{
			name:              "flags interspersed",
			input:             []string{"ls", "--dry", "-la", "--verbose"},
			expectedDry:       true,
			expectedVerbose:   true,
			expectedRemaining: []string{"ls", "-la"},
		},
		{
			name:              "empty input",
			input:             []string{},
			expectedDry:       false,
			expectedVerbose:   false,
			expectedRemaining: []string{},
		},
		{
			name:              "only global flags",
			input:             []string{"--dry", "--verbose"},
			expectedDry:       true,
			expectedVerbose:   true,
			expectedRemaining: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flags, remaining := ExtractGlobalFlags(tc.input)

			if flags.DryRun != tc.expectedDry {
				t.Errorf("Expected DryRun=%v, got %v", tc.expectedDry, flags.DryRun)
			}
			if flags.Verbose != tc.expectedVerbose {
				t.Errorf("Expected Verbose=%v, got %v", tc.expectedVerbose, flags.Verbose)
			}
			if !sliceEqual(remaining, tc.expectedRemaining) {
				t.Errorf("Expected remaining %v, got %v", tc.expectedRemaining, remaining)
			}
		})
	}
}
