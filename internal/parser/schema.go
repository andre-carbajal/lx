package parser

type CommandSchema struct {
	Args      ArgsSchema
	AllowPipe bool
	MaxArgs   int
	MaxFlags  int
}

type ArgsSchema struct {
	Type     string
	Multiple bool
	Optional bool
	MinCount int
}

func DefaultSchema() *CommandSchema {
	return &CommandSchema{
		Args: ArgsSchema{
			Type:     "free",
			Multiple: true,
			Optional: true,
		},
		AllowPipe: true,
		MaxArgs:   0,
		MaxFlags:  0,
	}
}

func ValidateAgainstSchema(cmd *ParsedCommand, schema *CommandSchema) []string {
	var errors []string

	if !schema.Args.Optional && len(cmd.Args) < schema.Args.MinCount {
		errors = append(errors,
			"error: insufficient arguments",
		)
	}

	if schema.MaxArgs > 0 && len(cmd.Args) > schema.MaxArgs {
		errors = append(errors,
			"error: too many arguments",
		)
	}

	if schema.MaxFlags > 0 && len(cmd.Flags) > schema.MaxFlags {
		errors = append(errors,
			"error: too many flags",
		)
	}

	if !schema.AllowPipe && cmd.Pipe != nil {
		errors = append(errors,
			"error: pipes not allowed for this command",
		)
	}

	return errors
}
