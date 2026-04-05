package translator

import (
	"strings"

	"github.com/andre-carbajal/lx/internal/parser"
	"github.com/andre-carbajal/lx/pkg/dict"
)

func renderTemplate(tpl string, cmd *parser.ParsedCommand, entry *dict.CommandEntry) (string, error) {
	result := tpl

	vars := map[string]string{
		"flags":   buildFlagsString(cmd, entry),
		"args":    strings.Join(cmd.Args, " "),
		"cmd":     cmd.Name,
		"pattern": getArgAt(cmd, 0),
		"files":   strings.Join(getArgsFrom(cmd, 1), " "),
		"source":  getArgAt(cmd, 0),
		"dest":    getArgAt(cmd, 1),
		"url":     getArgAt(cmd, 0),
		"pid":     getArgAt(cmd, 0),
		"count":   getArgAt(cmd, 0),
		"offset":  "1",
	}

	for key, value := range vars {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result, nil
}

func buildFlagsString(cmd *parser.ParsedCommand, entry *dict.CommandEntry) string {
	if len(cmd.Flags) == 0 {
		return ""
	}

	var flagParts []string

	for _, flag := range cmd.Flags {
		translation, ok := entry.Flags[flag]
		if !ok {
			flagParts = append(flagParts, flag)
			continue
		}

		flagParts = append(flagParts, translation.PS)
	}

	var result []string
	for _, f := range flagParts {
		if f != "" {
			result = append(result, f)
		}
	}

	return strings.Join(result, " ")
}

func getArgAt(cmd *parser.ParsedCommand, index int) string {
	if index < 0 || index >= len(cmd.Args) {
		return ""
	}
	return cmd.Args[index]
}

func getArgsFrom(cmd *parser.ParsedCommand, fromIndex int) []string {
	if fromIndex >= len(cmd.Args) {
		return []string{}
	}
	return cmd.Args[fromIndex:]
}
