package handler

import (
	"fmt"
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

// Help returns a brief help string
func (h *Gif) Help(short bool) string {
	if short {
		return "Returns the first giphy result for the supplied search term"
	} else {
		return fmt.Sprintf("Usage: `%vgif Jason Momoa`\n\nReturn: Probably a pretty hot gif", celexacreams.Prefix)
	}
}


// Handle shows the first gif returned by the supplied search string
func (h *Gif) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	if len(command) <= 1 {
		return "You should supply a search string, what do you think I am, a mind reader " + m.Author.Mention() + "?", "", make([]byte, 0), nil
	}
	searchString := strings.Join(command[1:], " ")
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
