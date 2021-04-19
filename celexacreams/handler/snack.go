package handler

import (
	"fmt"
	"github.com/morganonbass/celexacreams/celexacreams"

	"github.com/bwmarrin/discordgo"
)

// Snack responds to "snack"
type Snack struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *Snack) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Snack) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Snack) Help(short bool) string {
	if short {
		return "_Feed me_"
	} else {
		return fmt.Sprintf("Usage: `%vsnack`\n\nReturn: A happy bot", celexacreams.Prefix)
	}
}


// Handle feeds the bot
func (h *Snack) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	url, err := celexacreams.GetRandomGIF("cat eating")
	if err != nil {
		return "", "", make([]byte, 0), &celexacreams.CelexaError{
			"GIF error: " + err.Error(),
		}
	}
	return url, "", make([]byte, 0), nil
}
