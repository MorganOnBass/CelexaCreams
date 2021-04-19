package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"net/url"
	"strings"
)

// Mc responds to "mc"
type Mc struct {
	R, D bool
}

// Should I reply to the invoking message?
func (h *Mc) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Mc) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Mc) Help(short bool) string {
	if short {
		return "Generate a minecraft achievement"
	} else {
		return fmt.Sprintf("Usage: `%vmc [@mention] achievement`\n\nReturn: A minecraft achievement. If invoked " +
			"with a mention, it will be awarded to the mentioned user, else to the invoker.",
			celexacreams.Prefix)
	}
}


// Handle returns a minecraft achievement with specified name to the caller to the user mentioned in the invoking message
func (h *Mc) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error, ) {
	args, _ := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if len(args) < 2 {
		return "", "", make([]byte, 0), &celexacreams.MissingArgsError{Message: args[0]}
	}
	var u string
	if len(m.Mentions) > 0 {
		u = m.Mentions[0].Username
	} else {
		u = m.Author.Username
	}
	args = celexacreams.RemoveString(args, fmt.Sprintf("@%s", u))
	u = url.QueryEscape(u)
	txt := url.QueryEscape(strings.Join(args[1:], " "))
	url := fmt.Sprintf("https://mcgen.herokuapp.com/a.php?i=1&h=Achievement-%s&t=%s", u, txt)
	img, err := celexacreams.DownloadImage(url)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	return "", "achievement_unlocked.png", img, nil
}
