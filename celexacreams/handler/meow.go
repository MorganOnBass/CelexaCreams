package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
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

// Help returns a brief help string
func (h *Meow) Help(short bool) string {
	if short {
		return "A prototype"
	} else {
		return fmt.Sprintf("Usage: `%vmeow`\n\nI'm a lazy dev and this command serves as a starting point for " +
			"implementing new commands so I don't have to write all the boilerplate every time. :woman_shrugging:",
			celexacreams.Prefix)
	}
}


// Handle meows back
func (h *Meow) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	return "_meeeeow_", "", make([]byte, 0), nil
}
