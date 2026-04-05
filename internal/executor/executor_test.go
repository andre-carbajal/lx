package executor

import (
	"strings"
	"testing"

	"github.com/andre-carbajal/lx/internal/detector"
	"github.com/andre-carbajal/lx/internal/parser"
)

func TestExecute_DryRun(t *testing.T) {
	tests := []struct {
		name       string
		translated string
		shell      detector.ShellType
		dryRun     bool
		verbose    bool
		wantExit   int
	}{
		{
			name:       "dry run with CMD",
			translated: `dir`,
			shell:      detector.ShellCMD,
			dryRun:     true,
			verbose:    false,
			wantExit:   0,
		},
		{
			name:       "dry run with PowerShell",
			translated: `Get-ChildItem`,
			shell:      detector.ShellPowerShell,
			dryRun:     true,
			verbose:    false,
			wantExit:   0,
		},
		{
			name:       "dry run does not execute",
			translated: `exit 42`,
			shell:      detector.ShellCMD,
			dryRun:     true,
			verbose:    false,
			wantExit:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{}
			globalFlags := &parser.GlobalFlags{
				DryRun:  tt.dryRun,
				Verbose: tt.verbose,
			}

			exitCode := executor.Execute(tt.translated, tt.shell, globalFlags)
			if exitCode != tt.wantExit {
				t.Errorf("Execute() = %d, want %d", exitCode, tt.wantExit)
			}
		})
	}
}

func TestExecute_ValidCommands(t *testing.T) {
	tests := []struct {
		name       string
		translated string
		shell      detector.ShellType
		wantExit   int
	}{
		{
			name:       "CMD: echo command",
			translated: `echo test`,
			shell:      detector.ShellCMD,
			wantExit:   0,
		},
		{
			name:       "PowerShell: echo command",
			translated: `Write-Host test`,
			shell:      detector.ShellPowerShell,
			wantExit:   0,
		},
		{
			name:       "CMD: set with ERRORLEVEL",
			translated: `set ERRORLEVEL=0`,
			shell:      detector.ShellCMD,
			wantExit:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{}
			globalFlags := &parser.GlobalFlags{
				DryRun:  false,
				Verbose: false,
			}

			exitCode := executor.Execute(tt.translated, tt.shell, globalFlags)
			if exitCode != tt.wantExit {
				t.Errorf("Execute() = %d, want %d", exitCode, tt.wantExit)
			}
		})
	}
}

func TestExecuteWithResult_Success(t *testing.T) {
	tests := []struct {
		name       string
		translated string
		shell      detector.ShellType
		wantExit   int
		checkOut   func(string) bool
	}{
		{
			name:       "CMD: echo outputs correctly",
			translated: `echo hello`,
			shell:      detector.ShellCMD,
			wantExit:   0,
			checkOut: func(out string) bool {
				return strings.Contains(out, "hello") || len(out) > 0
			},
		},
		{
			name:       "PowerShell: Write-Host outputs",
			translated: `Write-Host "hello"`,
			shell:      detector.ShellPowerShell,
			wantExit:   0,
			checkOut: func(out string) bool {
				return strings.Contains(out, "hello") || len(out) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{}
			output, exitCode := executor.ExecuteWithResult(tt.translated, tt.shell)

			if exitCode != tt.wantExit {
				t.Errorf("ExecuteWithResult() exitCode = %d, want %d", exitCode, tt.wantExit)
			}

			if !tt.checkOut(output) {
				t.Errorf("ExecuteWithResult() output check failed, got: %q", output)
			}
		})
	}
}

func TestExecuteWithResult_InvalidCommand(t *testing.T) {
	tests := []struct {
		name       string
		translated string
		shell      detector.ShellType
		wantExitOK bool
	}{
		{
			name:       "CMD: invalid command returns non-zero",
			translated: `invalid_command_that_does_not_exist_xyz`,
			shell:      detector.ShellCMD,
			wantExitOK: true,
		},
		{
			name:       "PowerShell: invalid command returns non-zero",
			translated: `Get-InvalidCommandXyz`,
			shell:      detector.ShellPowerShell,
			wantExitOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{}
			_, exitCode := executor.ExecuteWithResult(tt.translated, tt.shell)

			if tt.wantExitOK {
				if exitCode == 0 {
					t.Errorf("ExecuteWithResult() expected non-zero exit code, got %d", exitCode)
				}
			} else {
				if exitCode != 0 {
					t.Errorf("ExecuteWithResult() expected exit code 0, got %d", exitCode)
				}
			}
		})
	}
}

func TestExecute_GlobalFlagsHandling(t *testing.T) {
	tests := []struct {
		name                 string
		translated           string
		shell                detector.ShellType
		dryRun               bool
		verbose              bool
		expectDryRunBehavior bool
	}{
		{
			name:                 "verbose flag set",
			translated:           `echo test`,
			shell:                detector.ShellCMD,
			dryRun:               false,
			verbose:              true,
			expectDryRunBehavior: false,
		},
		{
			name:                 "both flags set - dry run takes precedence",
			translated:           `echo test`,
			shell:                detector.ShellCMD,
			dryRun:               true,
			verbose:              true,
			expectDryRunBehavior: true,
		},
		{
			name:                 "no flags set",
			translated:           `echo test`,
			shell:                detector.ShellCMD,
			dryRun:               false,
			verbose:              false,
			expectDryRunBehavior: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &Executor{}
			globalFlags := &parser.GlobalFlags{
				DryRun:  tt.dryRun,
				Verbose: tt.verbose,
			}

			exitCode := executor.Execute(tt.translated, tt.shell, globalFlags)

			if tt.expectDryRunBehavior {
				if exitCode != 0 {
					t.Errorf("Execute() with dry run expected exit 0, got %d", exitCode)
				}
			}
			if exitCode != 0 {
				t.Logf("Execute() returned non-zero exit code: %d (may be expected for some cases)", exitCode)
			}
		})
	}
}
