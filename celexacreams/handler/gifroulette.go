package handler

import (
	"strings"

	"github.com/morganonbass/celexacreams/celexacreams"

	"github.com/bwmarrin/discordgo"
)

// GifRoulette responds to "GifRoulette"
type GifRoulette struct{}

// Handle shows a random gif returned by the supplied search string
func (h *GifRoulette) Handle(m *discordgo.MessageCreate) (string, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if len(command) <= 1 {
		return "You should supply a search string, what do you think I am, a mind reader " + m.Author.Mention() + "?", nil
	}
	searchString := strings.Join(command[1:], " ")
	if err != nil {
		return "", err
	}
	url, err := celexacreams.GetRandomGIF(searchString)
	if err != nil {
		return "", &celexacreams.CelexaError{
			"GIF error: " + err.Error(),
		}
	}
	return m.Author.Mention() + " " + url, nil
}
