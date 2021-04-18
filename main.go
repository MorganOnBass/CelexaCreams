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
	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	discordCreateMessageHandled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "discord_createmessage_handled",
		Help: "The total number of CreateMessage callbacks handled",
	})

	celexacreamsHandlers = map[string]celexacreams.Handler{
		// Handlers are structs, properties R and D control whether bot replies and deletes invocation respectively.
		// Setting both to true would be an error.
		"snack":       &handler.Snack{R: true},       // Feeds the bot
		"meow":        &handler.Meow{R: true},        // A prototype
		"gif":         &handler.Gif{R: true},         // Return the first giphy result for a search string
		"gifroulette": &handler.GifRoulette{R: true}, // return a random giphy result for a search string
		"magik":       &handler.Magik{R: true},       // magik an image, optional numeric argument specifies how much magik
		"deepfry":     &handler.DeepFry{R: true},     // Deep fry an image
		"haah":        &handler.Haah{R: true},        // Crop an image left half and mirror about Y axis
		"hooh":        &handler.Hooh{R: true},        // Crop an image top half and mirror about X axis
		"ahha":        &handler.Ahha{R: true},        // Crop an image right half and mirror about Y axis
		"ohho":        &handler.Ohho{R: true},        // Crop an image bottom half and mirror about X axis
		"aesthetic":   &handler.Aesthetic{D: true},   // Take some text and make it ａ ｅ ｓ ｔ ｈ ｅ ｔ ｉ ｃ
		"aesthetics":  &handler.Aesthetic{D: true},   // An alias because people keep trying to invoke it like this
		"jpeg":        &handler.Jpeg{R: true},        // Adds more jpeg
		"mc":          &handler.Mc{D: true},          // Builds a minecraft achievement
		"spin":        &handler.Spin{R: true},        // Crops to a circle and returns a spinning gif
	}
)

var Prefix string

func serveMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func main() {
	go serveMetrics()

	p, ok := os.LookupEnv("PREFIX")
	if !ok {
		log.Fatal("PREFIX is not set")
	}
	Prefix = p
	celexacreams.Prefix = Prefix

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
	imagick.Initialize()
	defer imagick.Terminate()

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

	if strings.HasPrefix(m.ContentWithMentionsReplaced(), Prefix) {
		if strings.HasPrefix(m.ContentWithMentionsReplaced(), Prefix+Prefix) {
			// now we don't invoke the bot when starting a message with ellipses
			return
		}
		rootHandler := &handler.Root{
			Handlers: celexacreamsHandlers,
		}

		response, err := rootHandler.Handle(m, c, s)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("failed to handle command: " + err.Error())
			response.Content = err.Error()
			response.Reference = &discordgo.MessageReference{
				ChannelID: c.ID,
				MessageID: m.ID,
				GuildID:   m.GuildID,
			}
		}

		_, err = s.ChannelMessageSendComplex(c.ID, response)
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
