package celexacreams

import (
	"github.com/bwmarrin/discordgo"
)

// Handler handles CelexaCreams commands
type Handler interface {
	Handle(*discordgo.MessageCreate) (string, error)
}
