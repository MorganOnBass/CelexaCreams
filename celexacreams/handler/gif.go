package handler

import (
	"strings"

	"github.com/morganonbass/celexacreams/celexacreams"

	"github.com/bwmarrin/discordgo"
)

// Gif responds to "Gif"
type Gif struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *Gif) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Gif) DeleteInvocation() bool {
	return h.D
}

// Handle shows the first gif returned by the supplied search string
func (h *Gif) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if len(command) <= 1 {
		return "You should supply a search string, what do you think I am, a mind reader " + m.Author.Mention() + "?", "", make([]byte, 0), nil
	}
	searchString := strings.Join(command[1:], " ")
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	url, err := celexacreams.GetGIF(searchString, 0)
	if err != nil {
		return "", "", make([]byte, 0), &celexacreams.CelexaError{
			"GIF error: " + err.Error(),
		}
	}
	r := new(discordgo.Message)
	r.Content = m.Author.Mention() + " " + url
	return url, "", make([]byte, 0), nil
}
