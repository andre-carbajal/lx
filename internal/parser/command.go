package parser

type ParsedCommand struct {
	Raw    string
	Name   string
	Flags  []string
	Args   []string
	Pipe   *ParsedCommand
	Schema *CommandSchema
}
