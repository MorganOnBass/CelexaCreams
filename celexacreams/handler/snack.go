package handler

import (
	"github.com/morganonbass/celexacreams/celexacreams"

	"github.com/bwmarrin/discordgo"
)

// Snack responds to "snack"
type Snack struct{}

// Handle shows snack
func (h *Snack) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, []byte, error) {
	url, err := celexacreams.GetRandomGIF("cat eating")
	if err != nil {
		return "", make([]byte, 0), &celexacreams.CelexaError{
			"GIF error: " + err.Error(),
		}
	}
	return url, make([]byte, 0), nil
}
