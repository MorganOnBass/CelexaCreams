package handler

import (
	"github.com/bwmarrin/discordgo"
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

// Handle response back
func (h *UwU) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (
	string,
	string,
	[]byte,
	error,
) {
	return "_What's this?_", "", make([]byte, 0), nil
}
