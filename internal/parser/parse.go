package parser

import (
	"errors"
	"strings"
)

func Parse(input string) (*ParsedCommand, error) {
	// Trim y validar input vacío
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("error: empty input")
	}

	// Detectar y parsear pipes (máximo 4 segmentos)
	segments := splitByPipe(input)
	if len(segments) > 4 {
		return nil, errors.New("error: too many pipes (max 4)")
	}

	// Parsear primer segmento
	cmd, err := parseSegment(segments[0])
	if err != nil {
		return nil, err
	}

	// Parsear y encadenar pipes
	current := cmd
	for i := 1; i < len(segments); i++ {
		next, err := parseSegment(segments[i])
		if err != nil {
			return nil, err
		}
		current.Pipe = next
		current = next
	}

	return cmd, nil
}

// splitByPipe divide por "|" no escapados
func splitByPipe(input string) []string {
	var segments []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for i, ch := range input {
		if !inQuote && (ch == '"' || ch == '\'') {
			inQuote = true
			quoteChar = ch
			current.WriteRune(ch)
		} else if inQuote && ch == quoteChar {
			// Revisar que no esté escapado
			if i > 0 && input[i-1] != '\\' {
				inQuote = false
			}
			current.WriteRune(ch)
		} else if !inQuote && ch == '|' {
			if current.Len() > 0 {
				segments = append(segments, strings.TrimSpace(current.String()))
				current.Reset()
			}
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		segments = append(segments, strings.TrimSpace(current.String()))
	}

	return segments
}

// parseSegment parsea un segmento individual
func parseSegment(segment string) (*ParsedCommand, error) {
	tokens := tokenize(segment)
	if len(tokens) == 0 {
		return nil, errors.New("error: empty segment")
	}

	cmd := &ParsedCommand{
		Raw:   segment,
		Name:  tokens[0],
		Flags: []string{},
		Args:  []string{},
	}

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if strings.HasPrefix(token, "-") && !isNumericFlag(token) {
			// Es un flag
			if strings.HasPrefix(token, "--") {
				// Flag larga
				cmd.Flags = append(cmd.Flags, token)
			} else if strings.HasPrefix(token, "-") && len(token) > 1 && token[1] != '-' {
				// Flag(s) corta(s)
				if len(token) == 2 {
					// Flag simple: -l
					cmd.Flags = append(cmd.Flags, token)
				} else if !startsWithDigit(token[1:]) {
					// Flags combinadas: -la → -l, -a (pero no -5 de head -5)
					for _, ch := range token[1:] {
						cmd.Flags = append(cmd.Flags, "-"+string(ch))
					}
				} else {
					// Es un argumento numérico como -5
					cmd.Args = append(cmd.Args, token)
				}
			}
		} else {
			// Es un argumento
			cmd.Args = append(cmd.Args, token)
		}
	}

	return cmd, nil
}

// isNumericFlag detecta si es un flag numérico como -5
func isNumericFlag(token string) bool {
	if !strings.HasPrefix(token, "-") || len(token) < 2 {
		return false
	}
	return startsWithDigit(token[1:])
}

// startsWithDigit verifica si la cadena comienza con un dígito
func startsWithDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= '0' && s[0] <= '9'
}

// tokenize divide el input en tokens respetando comillas
func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for i, ch := range input {
		if !inQuote && (ch == '"' || ch == '\'') {
			inQuote = true
			quoteChar = ch
			// NO incluir la comilla en el token
		} else if inQuote && ch == quoteChar {
			// NO incluir la comilla en el token
			if i > 0 && input[i-1] != '\\' {
				inQuote = false
			}
		} else if !inQuote && (ch == ' ' || ch == '\t') {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
