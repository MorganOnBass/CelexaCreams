package handler

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	celexaCreamsCommandHandledSuccess = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "celexaCreams_command_handled_success",
			Help: "The total number of CelexaCreams commands handled successfully",
		},
		[]string{"command"},
	)
	celexaCreamsCommandHandledError = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "celexaCreams_command_handled_error",
			Help: "The total number of CelexaCreams commands handled which errored",
		},
		[]string{"command"},
	)
)

// Root is the handler to rule all handlers
type Root struct {
	Handlers map[string]celexacreams.Handler
}

// Handle determines the command handler
func (h *Root) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (*discordgo.MessageSend, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if err != nil {
		return &discordgo.MessageSend{}, err
	}
	log.Info(command)

	handler, ok := h.Handlers[command[0]]
	if !ok {
		return &discordgo.MessageSend{}, &celexacreams.CommandNotFoundError{command[0]}
	}
	if handler.DeleteInvocation() {
		defer s.ChannelMessageDelete(c.ID, m.ID)
	}

	response, filename, pic, err := handler.Handle(m, c, s)
	if err != nil {
		celexaCreamsCommandHandledError.WithLabelValues(command[0]).Inc()
	}
	celexaCreamsCommandHandledSuccess.WithLabelValues(command[0]).Inc()

	ref := discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: c.ID,
		GuildID:   m.GuildID,
	}

	r := discordgo.MessageSend{
		Content:   response,
	}

	if handler.Reply() {
		r.Reference = &ref
	}

	if len(pic) > 0 {
		file := discordgo.File{
			Reader: bytes.NewReader(pic),
			Name:   filename,
		}
		r.Files = []*discordgo.File{&file}
	}

	return &r, err
}
