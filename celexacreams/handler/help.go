package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
)

// Help responds to "help"
type Help struct {
	R, D bool
}

// Should I reply to the invoking message?
func (h *Help) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Help) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Help) Help(short bool) string {
	if short {
		return "Get help about a specific command"
	} else {
		return fmt.Sprintf("Usage: `%vhelp command`\n\nReturn: Usage information for the requested command. Use %v" +
			"commands for a list of currently implemented commands.",
			celexacreams.Prefix, celexacreams.Prefix)
	}
}


// Handle returns help for the requested command
func (h *Help) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	if len(command) <= 1 || len(command) > 2 {
		return h.Help(false), "", make([]byte, 0), nil
	}
	if _, ok := celexacreams.Commands[command[1]]; !ok {
		return fmt.Sprintf("%v is not a command. Pull requests are welcome, though. :)", command[1]), "", make([]byte, 0), nil
	}
	return celexacreams.Commands[command[1]].Help(false), "", make([]byte, 0), nil
}
