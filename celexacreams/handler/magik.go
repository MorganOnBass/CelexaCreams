package handler

import (
	"github.com/morganonbass/celexacreams/celexacreams"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// time="2021-03-26T00:17:13Z" level=info msg="[magik]"
// time="2021-03-26T00:17:13Z" level=info msg="Debug noise!" m.Attachments="[]" m.Content=.magik m.ID=824799114511253554
// time="2021-03-26T00:17:53Z" level=info msg="[magik]"
// time="2021-03-26T00:17:53Z" level=info msg="Debug noise!" m.Attachments="[0xc00005a3c0]" m.Content=.magik m.ID=824799285496381500
// time="2021-03-26T00:24:33Z" level=info msg="[magik https://www.ikea.com/gb/en/images/products/blahaj-soft-toy-shark__0710175_pe727378_s5.jpg?f=xxl]"
// time="2021-03-26T00:24:33Z" level=info msg="Debug noise!" m.Attachments="[]" m.Content=".magik https://www.ikea.com/gb/en/images/products/blahaj-soft-toy-shark__0710175_pe727378_s5.jpg?f=xxl" m.ID=824800960622231583

// Magik responds to "magik"
type Magik struct{}

// Handle meows back
func (h *Magik) Handle(m *discordgo.MessageCreate) (string, error) {
	image, err := celexacreams.FindNearestImage(m)
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"m.ID":          m.ID,
		"m.Attachments": m.Attachments,
		"m.Content":     m.ContentWithMentionsReplaced(),
		"image":         image,
	}).Info("Debug noise!")

	return m.Author.Mention() + " Not implemented yet, but probably logging to console for Morgan to debug", nil
}
