package parser

import (
	"testing"
)

func TestValidateAgainstSchema(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *ParsedCommand
		schema      *CommandSchema
		shouldError bool
		errorCount  int
	}{
		{
			name: "valid: optional args with zero args",
			cmd: &ParsedCommand{
				Name:  "ls",
				Flags: []string{"-l"},
				Args:  []string{},
			},
			schema: &CommandSchema{
				Args: ArgsSchema{
					Type:     "paths",
					Optional: true,
				},
				AllowPipe: true,
			},
			shouldError: false,
			errorCount:  0,
		},
		{
			name: "valid: required args with args",
			cmd: &ParsedCommand{
				Name:  "grep",
				Flags: []string{},
				Args:  []string{"pattern", "file.txt"},
			},
			schema: &CommandSchema{
				Args: ArgsSchema{
					Type:     "pattern+files",
					Optional: false,
					MinCount: 1,
				},
				AllowPipe: true,
			},
			shouldError: false,
			errorCount:  0,
		},
		{
			name: "invalid: required args with zero args",
			cmd: &ParsedCommand{
				Name:  "grep",
				Flags: []string{},
				Args:  []string{},
			},
			schema: &CommandSchema{
				Args: ArgsSchema{
					Type:     "pattern+files",
					Optional: false,
					MinCount: 1,
				},
				AllowPipe: true,
			},
			shouldError: true,
			errorCount:  1,
		},
		{
			name: "invalid: max args exceeded",
			cmd: &ParsedCommand{
				Name:  "touch",
				Flags: []string{},
				Args:  []string{"file1.txt", "file2.txt", "file3.txt"},
			},
			schema: &CommandSchema{
				Args: ArgsSchema{
					Type:     "paths",
					Optional: false,
					MinCount: 1,
				},
				MaxArgs: 2,
			},
			shouldError: true,
			errorCount:  1,
		},
		{
			name: "invalid: pipes not allowed",
			cmd: &ParsedCommand{
				Name:  "echo",
				Flags: []string{},
				Args:  []string{"hello"},
				Pipe: &ParsedCommand{
					Name: "wc",
				},
			},
			schema: &CommandSchema{
				Args: ArgsSchema{
					Type:     "free",
					Optional: true,
				},
				AllowPipe: false,
			},
			shouldError: true,
			errorCount:  1,
		},
		{
			name: "valid: pipes allowed",
			cmd: &ParsedCommand{
				Name:  "ls",
				Flags: []string{"-l"},
				Args:  []string{},
				Pipe: &ParsedCommand{
					Name: "grep",
					Args: []string{"test"},
				},
			},
			schema: &CommandSchema{
				Args: ArgsSchema{
					Type:     "paths",
					Optional: true,
				},
				AllowPipe: true,
			},
			shouldError: false,
			errorCount:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errors := ValidateAgainstSchema(tc.cmd, tc.schema)

			if tc.shouldError && len(errors) == 0 {
				t.Errorf("Expected errors, got none")
			}
			if !tc.shouldError && len(errors) > 0 {
				t.Errorf("Expected no errors, got %d: %v", len(errors), errors)
			}
			if len(errors) != tc.errorCount {
				t.Errorf("Expected %d errors, got %d: %v", tc.errorCount, len(errors), errors)
			}
		})
	}
}

func TestDefaultSchema(t *testing.T) {
	schema := DefaultSchema()

	if schema == nil {
		t.Fatal("DefaultSchema returned nil")
	}

	if !schema.Args.Optional {
		t.Error("DefaultSchema should have optional args")
	}

	if !schema.Args.Multiple {
		t.Error("DefaultSchema should allow multiple args")
	}

	if !schema.AllowPipe {
		t.Error("DefaultSchema should allow pipes")
	}
}
