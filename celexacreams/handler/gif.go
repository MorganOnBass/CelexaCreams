package handler

import (
	"strings"

	"github.com/morganonbass/celexacreams/celexacreams"

	"github.com/bwmarrin/discordgo"
)

// Gif responds to "Gif"
type Gif struct{}

// Handle shows snack
func (h *Gif) Handle(m *discordgo.MessageCreate) (string, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if len(command) <= 1 {
		return "You should supply a search string, what do you think I am, a mind reader " + m.Author.Mention() + "?", nil
	}
	searchString := strings.Join(command[1:], " ")
	if err != nil {
		return "", err
	}
	url, err := celexacreams.GetGIF(searchString)
	if err != nil {
		return "", &celexacreams.CelexaError{
			"GIF error: " + err.Error(),
		}
	}
	return m.Author.Mention() + " " + url, nil
}
