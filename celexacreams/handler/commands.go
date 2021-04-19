package handler

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"strings"
	"unicode/utf8"
)

// Commands responds to "commands"
type Commands struct {
	R, D bool
}

// Should I reply to the invoking message?
func (h *Commands) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Commands) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Commands) Help(short bool) string {
	if short {
		return "List commands, with a brief description"
	} else {
		return fmt.Sprintf("Usage: `%vcommands`\n\nReturn: A list of currently implemented commands.", celexacreams.Prefix)
	}
}


// Handle returns a list of commands
func (h *Commands) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("Commands:\n\n")
	for k, v := range celexacreams.Commands {
		pad := strings.Repeat(" ", 18 - utf8.RuneCountInString(k))
		buffer.WriteString(fmt.Sprintf("***%s:***%s%s\n", k, pad, v.Help(true)))
	}
	return buffer.String(), "", make([]byte, 0), nil
}
