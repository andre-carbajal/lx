package parser

type GlobalFlags struct {
	DryRun  bool
	Verbose bool
}

func ExtractGlobalFlags(args []string) (*GlobalFlags, []string) {
	flags := &GlobalFlags{}
	var remaining []string

	for _, arg := range args {
		switch arg {
		case "--dry":
			flags.DryRun = true
		case "--verbose", "-v":
			flags.Verbose = true
		default:
			remaining = append(remaining, arg)
		}
	}

	return flags, remaining
}
