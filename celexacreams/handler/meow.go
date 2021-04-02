package handler

import (
	"github.com/bwmarrin/discordgo"
)

// Meow responds to "meow"
type Meow struct {
	R, D bool
}

// Should I reply to the invoking message?
func (h *Meow) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Meow) DeleteInvocation() bool {
	return h.D
}

// Handle meows back
func (h *Meow) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (
	string,
	string,
	[]byte,
	error,
) {
	return "_meeeeow_", "", make([]byte, 0), nil
}
