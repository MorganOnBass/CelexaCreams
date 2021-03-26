package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"github.com/morganonbass/celexacreams/celexacreams/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	discordCreateMessageHandled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "discord_createmessage_handled",
		Help: "The total number of CreateMessage callbacks handled",
	})

	celexacreamsHandlers = map[string]celexacreams.Handler{
		"snack":       &handler.Snack{},
		"meow":        &handler.Meow{},
		"gif":         &handler.Gif{},
		"gifroulette": &handler.GifRoulette{},
		"magik":       &handler.Magik{},
	}
)

func serveMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func main() {
	go serveMetrics()

	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		log.Fatal("DISCORD_TOKEN not set")
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to initialise Discord bot")
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to open Discord session")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func mentionsCelexaCreams(m *discordgo.MessageCreate, id string) bool {
	for _, m := range m.Mentions {
		if m.ID == id {
			return true
		}
	}
	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	discordCreateMessageHandled.Inc()

	if m.Author.ID == s.State.User.ID {
		return
	}

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	if strings.HasPrefix(m.ContentWithMentionsReplaced(), ".") {
		rootHandler := &handler.Root{
			Handlers: celexacreamsHandlers,
		}

		response, err := rootHandler.Handle(m)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("failed to handle command: " + err.Error())
			response = err.Error()
		}

		_, err = s.ChannelMessageSend(c.ID, response)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("failed to send message")
		}
		return
	}

	if mentionsCelexaCreams(m, s.State.User.ID) {
		err := s.MessageReactionAdd(c.ID, m.ID, "😻")
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("failed to add reaction")
		}
	}
}
