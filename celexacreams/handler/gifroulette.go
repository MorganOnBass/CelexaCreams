package handler

import (
	"fmt"
	"strings"

	"github.com/morganonbass/celexacreams/celexacreams"

	"github.com/bwmarrin/discordgo"
)

// GifRoulette responds to "GifRoulette"
type GifRoulette struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *GifRoulette) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *GifRoulette) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *GifRoulette) Help(short bool) string {
	if short {
		return "Returns a random giphy result for the supplied search term"
	} else {
		return fmt.Sprintf("Usage: `%vgifroulette Jason Momoa`\n\nReturn: Probably a pretty hot gif", celexacreams.Prefix)
	}
}


// Handle shows a random gif returned by the supplied search string
func (h *GifRoulette) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if len(command) <= 1 {
		return "You should supply a search string, what do you think I am, a mind reader " + m.Author.Mention() + "?", "", make([]byte, 0), nil
	}
	searchString := strings.Join(command[1:], " ")
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	url, err := celexacreams.GetRandomGIF(searchString)
	if err != nil {
		return "", "", make([]byte, 0), &celexacreams.CelexaError{
			"GIF error: " + err.Error(),
		}
	}
	return url, "", make([]byte, 0), nil
}
