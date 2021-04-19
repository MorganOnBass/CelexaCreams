package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"strings"
)

// Aesthetic responds to "aesthetic"
type Aesthetic struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *Aesthetic) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Aesthetic) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Aesthetic) Help(short bool) string {
	if short {
		return "ａ ｅ ｓ ｔ ｈ ｅ ｔ ｉ ｃ"
	} else {
		return fmt.Sprintf("Usage: `%vaesthetic aesthetic`\n\nReturn: ａ ｅ ｓ ｔ ｈ ｅ ｔ ｉ ｃ", celexacreams.Prefix)
	}
}

// Handle returns the input text but ａ ｅ ｓ ｔ ｈ ｅ ｔ ｉ ｃ
func (h *Aesthetic) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (
	string,
	string,
	[]byte,
	error,
) {
	args, _ := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if len(args) < 2 {
		return "", "", make([]byte, 0), &celexacreams.MissingArgsError{Message: args[0]}
	}
	input := strings.Join(args[1:], " ")
	iRunes := []rune(input)
	var oRunes []rune
	for _, r := range iRunes {
		if r < 127 && r > 32 {
			oRunes = append(oRunes, r + 65248, 32)
		} else if r == 32 {
			oRunes = append(oRunes, r, r, r)
		} else {
			oRunes = append(oRunes, 32, r, 32)
		}
	}

	return string(oRunes), "", make([]byte, 0), nil
}
