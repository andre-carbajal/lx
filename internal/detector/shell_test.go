package detector

import (
	"errors"
	"testing"
)

type MockProcessProvider struct {
	parentName string
	err        error
}

func (m *MockProcessProvider) GetParentProcessName(pid int) (string, error) {
	return m.parentName, m.err
}

func TestDetectCMD(t *testing.T) {
	provider := &MockProcessProvider{parentName: "cmd.exe"}
	detector := New(provider)
	if got := detector.Detect(); got != ShellCMD {
		t.Errorf("Expected ShellCMD, got %v", got)
	}
}

func TestDetectPowerShell(t *testing.T) {
	provider := &MockProcessProvider{parentName: "powershell.exe"}
	detector := New(provider)
	if got := detector.Detect(); got != ShellPowerShell {
		t.Errorf("Expected ShellPowerShell, got %v", got)
	}
}

func TestDetectPwsh(t *testing.T) {
	provider := &MockProcessProvider{parentName: "pwsh.exe"}
	detector := New(provider)
	if got := detector.Detect(); got != ShellPowerShell {
		t.Errorf("Expected ShellPowerShell, got %v", got)
	}
}

func TestDetectUnknown(t *testing.T) {
	provider := &MockProcessProvider{parentName: "bash.exe"}
	detector := New(provider)
	if got := detector.Detect(); got != ShellUnknown {
		t.Errorf("Expected ShellUnknown, got %v", got)
	}
}

// Test para override via env var (se implementa con t.Setenv en Go 1.21+)
func TestDetectEnvOverride(t *testing.T) {
	t.Setenv("LX_SHELL", "cmd")
	provider := &MockProcessProvider{parentName: "bash.exe"}
	detector := New(provider)
	if got := detector.Detect(); got != ShellCMD {
		t.Errorf("Expected ShellCMD (override), got %v", got)
	}
}

func TestDetectEnvOverridePowerShell(t *testing.T) {
	t.Setenv("LX_SHELL", "powershell")
	provider := &MockProcessProvider{parentName: "cmd.exe"}
	detector := New(provider)
	if got := detector.Detect(); got != ShellPowerShell {
		t.Errorf("Expected ShellPowerShell (override), got %v", got)
	}
}

func TestDetectError(t *testing.T) {
	provider := &MockProcessProvider{err: errors.New("failed to get process name")}
	detector := New(provider)
	if got := detector.Detect(); got != ShellUnknown {
		t.Errorf("Expected ShellUnknown on error, got %v", got)
	}
}

func TestShellTypeString(t *testing.T) {
	tests := []struct {
		shellType ShellType
		expected  string
	}{
		{ShellCMD, "cmd"},
		{ShellPowerShell, "powershell"},
		{ShellUnknown, "unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			if got := tc.shellType.String(); got != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestWindowsProcessProvider_Creation(t *testing.T) {
	provider := NewWindowsProcessProvider()
	if provider == nil {
		t.Error("NewWindowsProcessProvider() returned nil")
	}
}

func TestDefaultProvider(t *testing.T) {
	provider := DefaultProvider()
	if provider == nil {
		t.Error("DefaultProvider() returned nil")
	}

	// Verify it's a ProcessProvider
	var _ ProcessProvider = provider
}
