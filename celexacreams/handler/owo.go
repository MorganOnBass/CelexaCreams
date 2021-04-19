package handler

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/morganonbass/celexacreams/celexacreams"
)

// OwO responds to "owo"
type OwO struct {
    R, D bool
}

// Should I reply to the invoking message?
func (h *OwO) Reply() bool {
    return h.R
}

// Should I delete the invoking message?
func (h *OwO) DeleteInvocation() bool {
    return h.D
}

// Help returns a brief help string
func (h *OwO) Help(short bool) string {
    if short {
        return "UwU"
    } else {
        return fmt.Sprintf("Usage: `%vowo`\n\nReturn: UwU", celexacreams.Prefix)
    }
}


// Handle response UwU
func (h *OwO) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (
    string,
    string,
    []byte,
    error,
) {
    return "UwU", "", make([]byte, 0), nil
}