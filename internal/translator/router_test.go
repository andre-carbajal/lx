package translator

import (
	"testing"

	"github.com/andre-carbajal/lx/internal/detector"
	"github.com/andre-carbajal/lx/internal/parser"
)

func TestNewRouter(t *testing.T) {
	router, err := New()
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	if router == nil {
		t.Fatal("Router is nil")
	}

	if router.dict == nil {
		t.Fatal("Router dictionary is nil")
	}
}

func TestTranslateLS(t *testing.T) {
	router, _ := New()

	tests := []struct {
		name     string
		input    string
		shell    detector.ShellType
		contains string
	}{
		{
			name:     "ls to CMD",
			input:    "ls",
			shell:    detector.ShellCMD,
			contains: "dir",
		},
		{
			name:     "ls to PowerShell",
			input:    "ls",
			shell:    detector.ShellPowerShell,
			contains: "Get-ChildItem",
		},
		{
			name:     "ls -l to PowerShell",
			input:    "ls -l",
			shell:    detector.ShellPowerShell,
			contains: "Get-ChildItem",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd, _ := parser.Parse(tc.input)
			result, err := router.Translate(cmd, tc.shell)

			if err != nil {
				t.Errorf("Translation failed: %v", err)
				return
			}

			if result.Source != SourceDict {
				t.Errorf("Expected SourceDict, got %v", result.Source)
			}

			if result.Translated == "" {
				t.Error("Translation is empty")
			}

			if !contains(result.Translated, tc.contains) {
				t.Errorf("Expected translation to contain %q, got %q", tc.contains, result.Translated)
			}
		})
	}
}

func TestTranslateUnknownCommand(t *testing.T) {
	router, _ := New()

	cmd, _ := parser.Parse("nonexistent")
	result, err := router.Translate(cmd, detector.ShellPowerShell)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.Source != SourceUnknown {
		t.Errorf("Expected SourceUnknown, got %v", result.Source)
	}

	if len(result.Warnings) == 0 {
		t.Error("Expected warnings for unknown command")
	}

	if result.Translated != "" {
		t.Errorf("Expected empty translation, got %q", result.Translated)
	}
}

func TestTranslateGREP(t *testing.T) {
	router, _ := New()

	cmd, _ := parser.Parse(`grep "error" file.txt`)
	result, err := router.Translate(cmd, detector.ShellPowerShell)

	if err != nil {
		t.Errorf("Translation failed: %v", err)
	}

	if result.Source != SourceDict {
		t.Errorf("Expected SourceDict, got %v", result.Source)
	}

	if !contains(result.Translated, "Select-String") {
		t.Errorf("Expected PowerShell Select-String, got %q", result.Translated)
	}

	if !contains(result.Translated, "error") {
		t.Errorf("Expected pattern in translation, got %q", result.Translated)
	}
}

func TestTranslateCAT(t *testing.T) {
	router, _ := New()

	cmd, _ := parser.Parse("cat /etc/passwd")
	result, err := router.Translate(cmd, detector.ShellCMD)

	if err != nil {
		t.Errorf("Translation failed: %v", err)
	}

	if result.Source != SourceDict {
		t.Errorf("Expected SourceDict, got %v", result.Source)
	}

	if !contains(result.Translated, "type") {
		t.Errorf("Expected CMD type command, got %q", result.Translated)
	}
}

func TestTranslateClear(t *testing.T) {
	router, _ := New()

	cmd, _ := parser.Parse("clear")
	resultCMD, _ := router.Translate(cmd, detector.ShellCMD)
	resultPS, _ := router.Translate(cmd, detector.ShellPowerShell)

	if !contains(resultCMD.Translated, "cls") {
		t.Errorf("Expected 'cls' in CMD translation, got %q", resultCMD.Translated)
	}

	if !contains(resultPS.Translated, "Clear-Host") {
		t.Errorf("Expected 'Clear-Host' in PS translation, got %q", resultPS.Translated)
	}
}

func TestTranslationSource(t *testing.T) {
	if SourceDict.String() != "dictionary" {
		t.Errorf("Expected 'dictionary', got %q", SourceDict.String())
	}

	if SourceUnknown.String() != "unknown" {
		t.Errorf("Expected 'unknown', got %q", SourceUnknown.String())
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) &&
		(s == substr || len(s) > 0 && (s[0:len(substr)] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
