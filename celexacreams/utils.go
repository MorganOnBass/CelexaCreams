package celexacreams

import (
	"fmt"
	"strings"
)

type CelexaError struct {
	Message string
}

func (e *CelexaError) Error() string {
	return fmt.Sprintf("celexafail: %s", e.Message)
}

type InvalidCommandFormatError struct {
	Message string
}

func (e *InvalidCommandFormatError) Error() string {
	return fmt.Sprintf("invalid command format: %s", e.Message)
}

type CommandNotFoundError struct {
	Message string
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("command not found: %s", e.Message)
}

func ExtractCommand(s string) ([]string, error) {
	splitInput := strings.Split(s, " ")

	for i, s := range splitInput {
		splitInput[i] = strings.TrimSpace(s)
	}

	if splitInput[0] != "@CelexaCreams" {
		return []string{}, &InvalidCommandFormatError{
			"first segment must be '@CelexaCreams'",
		}
	}

	if len(splitInput) < 2 {
		return []string{}, &InvalidCommandFormatError{
			"command must be at least 2 segments, @CelexaCreams <command>",
		}
	}

	return splitInput, nil
}
