package translator

import (
	"fmt"
	"time"

	"github.com/andre-carbajal/lx/internal/detector"
	"github.com/andre-carbajal/lx/internal/parser"
	"github.com/andre-carbajal/lx/pkg/dict"
)

type TranslationSource int

const (
	SourceDict TranslationSource = iota
	SourceUnknown
)

type TranslationResult struct {
	Translated string
	Source     TranslationSource
	Shell      detector.ShellType
	Warnings   []string
	Latency    time.Duration
}

type Router struct {
	dict *dict.Dictionary
}

func New() (*Router, error) {
	d, err := dict.Load()
	if err != nil {
		return nil, err
	}

	return &Router{
		dict: d,
	}, nil
}

func (r *Router) Translate(cmd *parser.ParsedCommand, shell detector.ShellType) (*TranslationResult, error) {
	start := time.Now()
	result := &TranslationResult{
		Shell:    shell,
		Warnings: []string{},
		Source:   SourceUnknown,
	}

	entry := r.dict.Get(cmd.Name)
	if entry == nil {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("command '%s' not found in dictionary", cmd.Name),
		)
		result.Latency = time.Since(start)
		return result, nil
	}

	var template string
	if shell == detector.ShellPowerShell {
		template = entry.PS
	} else {
		template = entry.CMD
	}

	translated, err := renderTemplate(template, cmd, entry)
	if err != nil {
		return nil, err
	}

	result.Translated = translated
	result.Source = SourceDict
	result.Latency = time.Since(start)

	return result, nil
}

func (ts TranslationSource) String() string {
	switch ts {
	case SourceDict:
		return "dictionary"
	case SourceUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}
