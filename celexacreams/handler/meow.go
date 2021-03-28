package handler

import (
	"github.com/bwmarrin/discordgo"
)

// Meow responds to "meow"
type Meow struct{}

// Handle meows back
func (h *Meow) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, []byte, error) {
	return "_meeeeow_", make([]byte, 0), nil
}
