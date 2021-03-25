package handler

import (
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
func (c *Root) Handle(m *discordgo.MessageCreate) (string, error) {
	command, err := celexacreams.ExtractCommand(m.ContentWithMentionsReplaced())
	if err != nil {
		return "", err
	}
	log.Info(command)

	handler, ok := c.Handlers[command[1]]
	if !ok {
		return "", &celexacreams.CommandNotFoundError{command[1]}
	}

	response, err := handler.Handle(m)
	if err != nil {
		celexaCreamsCommandHandledError.WithLabelValues(command[1]).Inc()
	}
	celexaCreamsCommandHandledSuccess.WithLabelValues(command[1]).Inc()

	return response, err
}
