package handler

import(
	"github.com/bwmarrin/discordgo"
)

// Meow responds to "meow"
type Meow struct {}

// Handle meows back
func (h *Meow) Handle(m *discordgo.MessageCreate) (string, error) {
	return m.Author.Mention() + " _meeeeow_", nil
}