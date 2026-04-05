package parser

import (
	"testing"
)

var testCases = []struct {
	name          string
	input         string
	expectedCmd   string
	expectedFlags []string
	expectedArgs  []string
	shouldError   bool
}{
	// Casos básicos
	{"cmd only", "ls", "ls", []string{}, []string{}, false},
	{"single flag short", "ls -l", "ls", []string{"-l"}, []string{}, false},
	{"combined flags", "ls -la", "ls", []string{"-l", "-a"}, []string{}, false},
	{"separate flags", "ls -l -a", "ls", []string{"-l", "-a"}, []string{}, false},
	{"long flag", "ls --all", "ls", []string{"--all"}, []string{}, false},

	// Argumentos
	{"single arg", "ls /home", "ls", []string{}, []string{"/home"}, false},
	{"flag and arg", "ls -l /home", "ls", []string{"-l"}, []string{"/home"}, false},
	{"multiple args", "ls -la /home /tmp", "ls", []string{"-l", "-a"}, []string{"/home", "/tmp"}, false},

	// Comillas
	{"double quotes", `grep "hello world" file.txt`, "grep", []string{}, []string{"hello world", "file.txt"}, false},
	{"single quotes", "grep 'hello world' file.txt", "grep", []string{}, []string{"hello world", "file.txt"}, false},

	// Flags como argumentos (find -name) - en realidad se parsean como flags
	// La decisión es: si empieza con "-", es un flag. El traductor sabrá qué hacer.
	{"flag as arg", `find /home -name "*.txt"`, "find", []string{"-n", "-a", "-m", "-e"}, []string{"/home", "*.txt"}, false},

	// Argumentos especiales
	{"sed regex", `sed 's/foo/bar/g' input.txt`, "sed", []string{}, []string{"s/foo/bar/g", "input.txt"}, false},
	{"awk braces", `awk '{print $1}' file.txt`, "awk", []string{}, []string{"{print $1}", "file.txt"}, false},

	// Espacios múltiples
	{"extra spaces", "  ls   -l  /home  ", "ls", []string{"-l"}, []string{"/home"}, false},

	// Variables de entorno (sin expansión)
	{"env var", "echo $HOME", "echo", []string{}, []string{"$HOME"}, false},

	// Numeric flag (head -5)
	{"numeric flag", "head -5", "head", []string{}, []string{"-5"}, false},

	// Rutas con espacios
	{"path with spaces double quotes", `ls "/home/user/my docs"`, "ls", []string{}, []string{"/home/user/my docs"}, false},

	// Errores
	{"empty input", "", "", nil, nil, true},
	// Nota: "ls /home/user/my docs" se parsea como ["ls", "/home/user/my", "docs"] - es válido
	// El usuario debería usar comillas si quiere un arg con espacios
}

func TestParse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := Parse(tc.input)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cmd.Name != tc.expectedCmd {
				t.Errorf("Expected name %q, got %q", tc.expectedCmd, cmd.Name)
			}

			if !sliceEqual(cmd.Flags, tc.expectedFlags) {
				t.Errorf("Expected flags %v, got %v", tc.expectedFlags, cmd.Flags)
			}

			if !sliceEqual(cmd.Args, tc.expectedArgs) {
				t.Errorf("Expected args %v, got %v", tc.expectedArgs, cmd.Args)
			}
		})
	}
}

func TestParseWithPipes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		segments int
	}{
		{"single pipe", "ls | grep txt", 2},
		{"double pipe", "cat file.txt | grep error | head -5", 3},
		{"pipe with complex args", `find . -name "*.txt" | xargs grep "hello" | wc -l`, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := Parse(tc.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Count pipe segments
			count := 1
			current := cmd
			for current.Pipe != nil {
				count++
				current = current.Pipe
			}

			if count != tc.segments {
				t.Errorf("Expected %d segments, got %d", tc.segments, count)
			}
		})
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
