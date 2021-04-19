package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
)

// UwU responds to "uwu"
type UwU struct {
	R, D bool
}

// Should I reply to the invoking message?
func (h *UwU) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *UwU) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *UwU) Help(short bool) string {
	if short {
		return "_What's this?_"
	} else {
		return fmt.Sprintf("Usage: `%vuwu`\n\nReturn: _What's this?_", celexacreams.Prefix)
	}
}


// Handle responds "_What's this?_"
func (h *UwU) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (
	string,
	string,
	[]byte,
	error,
) {
	return "_What's this?_", "", make([]byte, 0), nil
}
