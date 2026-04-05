package dict

import (
	"testing"
)

func TestLoadDictionary(t *testing.T) {
	dict, err := Load()
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	if dict == nil {
		t.Fatal("Dictionary is nil")
	}

	if dict.Count() != 26 {
		t.Errorf("Expected 26 commands, got %d", dict.Count())
	}
}

func TestDictionaryGet(t *testing.T) {
	dict, _ := Load()

	tests := []struct {
		name     string
		cmd      string
		exists   bool
		hasFlags bool
	}{
		{"ls exists", "ls", true, true},
		{"grep exists", "grep", true, true},
		{"cat exists", "cat", true, false},
		{"echo exists", "echo", true, true},
		{"unknown command", "nonexistent", false, false},
		{"chmod not in dict", "chmod", false, false},
		{"sudo not in dict", "sudo", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			entry := dict.Get(tc.cmd)

			if tc.exists && entry == nil {
				t.Errorf("Expected command %q to exist, but got nil", tc.cmd)
			}
			if !tc.exists && entry != nil {
				t.Errorf("Expected command %q to not exist, but got %v", tc.cmd, entry)
			}

			if tc.exists && tc.hasFlags && len(entry.Flags) == 0 {
				t.Errorf("Expected command %q to have flags, but got none", tc.cmd)
			}
		})
	}
}

func TestDictionaryExists(t *testing.T) {
	dict, _ := Load()

	if !dict.Exists("ls") {
		t.Error("Expected ls to exist")
	}

	if !dict.Exists("grep") {
		t.Error("Expected grep to exist")
	}

	if dict.Exists("nonexistent") {
		t.Error("Expected nonexistent to not exist")
	}
}

func TestCommandEntry(t *testing.T) {
	dict, _ := Load()

	ls := dict.Get("ls")
	if ls == nil {
		t.Fatal("ls command not found")
	}

	if ls.Description == "" {
		t.Error("ls description is empty")
	}
	if ls.CMD == "" {
		t.Error("ls CMD is empty")
	}
	if ls.PS == "" {
		t.Error("ls PS is empty")
	}

	if ls.Schema.Args.Type == "" {
		t.Error("ls schema args type is empty")
	}
	if !ls.Schema.AllowPipe {
		t.Error("ls should allow pipes")
	}
}

func TestFlagTranslations(t *testing.T) {
	dict, _ := Load()

	ls := dict.Get("ls")
	if ls == nil {
		t.Fatal("ls command not found")
	}

	lFlag, ok := ls.Flags["-l"]
	if !ok {
		t.Error("ls should have -l flag")
	}
	if lFlag.PS == "" {
		t.Error("-l PS translation is empty")
	}

	laFlag, ok := ls.Flags["-la"]
	if !ok {
		t.Error("ls should have -la flag")
	}
	if laFlag.CMD == "" || laFlag.PS == "" {
		t.Error("-la flag translations are incomplete")
	}
}

func TestAllCommandsHaveSchema(t *testing.T) {
	dict, _ := Load()

	for _, cmdName := range dict.ListCommands() {
		cmd := dict.Get(cmdName)
		if cmd == nil {
			t.Errorf("Command %q not found", cmdName)
			continue
		}

		if cmd.Schema.Args.Type == "" {
			t.Errorf("Command %q has empty schema args type", cmdName)
		}
	}
}

func TestListCommands(t *testing.T) {
	dict, _ := Load()

	commands := dict.ListCommands()
	if len(commands) != 26 {
		t.Errorf("Expected 26 commands, got %d", len(commands))
	}

	known := map[string]bool{
		"ls": false, "cd": false, "grep": false, "cat": false, "echo": false,
	}

	for _, cmd := range commands {
		if _, ok := known[cmd]; ok {
			known[cmd] = true
		}
	}

	for cmd, found := range known {
		if !found {
			t.Errorf("Expected command %q to be in list", cmd)
		}
	}
}
