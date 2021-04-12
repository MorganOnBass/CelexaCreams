package handler

import (
    "github.com/bwmarrin/discordgo"
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

// Handle response back
func (h *OwO) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (
    string,
    string,
    []byte,
    error,
) {
    return "UwU", "", make([]byte, 0), nil
}