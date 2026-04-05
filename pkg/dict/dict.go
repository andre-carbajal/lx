package dict

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed commands.yaml
var commandsYAML string

type CommandEntry struct {
	Description string                     `yaml:"description"`
	CMD         string                     `yaml:"cmd"`
	PS          string                     `yaml:"ps"`
	Schema      CommandSchemaYAML          `yaml:"schema"`
	Flags       map[string]FlagTranslation `yaml:"flags"`
}

type CommandSchemaYAML struct {
	Args      ArgsSchemaYAML `yaml:"args"`
	AllowPipe bool           `yaml:"allowPipe"`
	MaxArgs   int            `yaml:"maxArgs,omitempty"`
	MaxFlags  int            `yaml:"maxFlags,omitempty"`
}

type ArgsSchemaYAML struct {
	Type     string `yaml:"type"`
	Multiple bool   `yaml:"multiple"`
	Optional bool   `yaml:"optional"`
	MinCount int    `yaml:"minCount,omitempty"`
}

type FlagTranslation struct {
	CMD string `yaml:"cmd"`
	PS  string `yaml:"ps"`
}

type Dictionary struct {
	Commands map[string]*CommandEntry
}

func Load() (*Dictionary, error) {
	var data struct {
		Commands map[string]*CommandEntry `yaml:"commands"`
	}

	err := yaml.Unmarshal([]byte(commandsYAML), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commands.yaml: %w", err)
	}

	return &Dictionary{
		Commands: data.Commands,
	}, nil
}

func (d *Dictionary) Get(name string) *CommandEntry {
	return d.Commands[name]
}

func (d *Dictionary) Exists(name string) bool {
	_, ok := d.Commands[name]
	return ok
}

func (d *Dictionary) Count() int {
	return len(d.Commands)
}

func (d *Dictionary) ListCommands() []string {
	commands := make([]string, 0, len(d.Commands))
	for name := range d.Commands {
		commands = append(commands, name)
	}
	return commands
}
