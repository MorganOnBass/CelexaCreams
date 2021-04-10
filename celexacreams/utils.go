package celexacreams

import (
	"fmt"
	"net/url"
	"strings"
)

var Prefix string

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

type NotAnImageError struct {
	Message string
}

func (e *NotAnImageError) Error() string {
	return fmt.Sprintf("not an image: %s", e.Message)
}

type ImageNotFoundError struct {
	Message string
}

func (e *ImageNotFoundError) Error() string {
	return fmt.Sprintf("Image not found in message: %s", e.Message)
}

type MissingArgsError struct {
	Message string
}

func (e *MissingArgsError) Error() string {
	return fmt.Sprintf("%s command expects an argument", e.Message)
}

func ExtractCommand(s string) ([]string, error) {
	splitInput := strings.Split(s, " ")
	splitInput[0] = strings.Trim(splitInput[0], Prefix)

	for i, s := range splitInput {
		splitInput[i] = strings.TrimSpace(s)
	}

	if len(splitInput) < 1 {
		return []string{}, &InvalidCommandFormatError{
			"There's no command here?",
		}
	}

	return splitInput, nil
}

func IsURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func RemoveString(s []string, v string) []string {
	i := FindString(s, v)
	if i != -1 {
		return append(s[:i], s[i+1:]...)
	} else {
		return s
	}
}

func FindString(s []string, v string) int {
	for i, k := range s {
		if k == v {
			return i
		}
	}
	return -1
}
