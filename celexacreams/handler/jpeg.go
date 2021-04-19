package handler

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	"github.com/morganonbass/celexacreams/celexacreams"
)

// Jpeg responds to "jpeg"
type Jpeg struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *Jpeg) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Jpeg) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Jpeg) Help(short bool) string {
	if short {
		return "Needs moar jpeg. Optional integer argument sets just how much jpeg to add"
	} else {
		return fmt.Sprintf("Usage: `%vjpeg [1-10]`\n\nYou may attach an image or link to an image to the invoking post," +
			" or reply to a post containing an image. If you do not, jpeg will process the most recent image in the channel.",
			celexacreams.Prefix)
	}
}


// Handle returns an image with moar jpeg. An optional integer argument specifies just how much more jpeg to add
func (h *Jpeg) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	args, _ := celexacreams.ExtractCommand(m.Content)
	var sauce int
	if len(args) > 1 {
		arg, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			arg = 10 // probably got '.jpeg <url>'
		}
		sauce = int(arg)
	} else {
		sauce = 10
	}
	if sauce < 1 || sauce > 10 {
		return "jpeg only accepts a number between 1 and 10", "", make([]byte, 0), nil
	} else {
		sauce = 11 - sauce // invert this because higher args to jpeg conversion make less jpeg
	}

	sTime := time.Now()
	URL, err := celexacreams.FindNearestImageURL(m, c, s)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	img, err := celexacreams.DownloadImage(URL)
	if err != nil {
		return "", "", make([]byte, 0), err
	}

	ref := discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: c.ID,
		GuildID:   m.GuildID,
	}

	r := discordgo.MessageSend{
		Content:   "Processing...",
		Reference: &ref,
	}
	msg, err := s.ChannelMessageSendComplex(c.ID, &r)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	defer s.ChannelMessageDelete(c.ID, msg.ID)

	i, err := imaging.Decode(bytes.NewReader(img), imaging.AutoOrientation(true))
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, i, imaging.JPEG, imaging.JPEGQuality(sauce))
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	output := buf.Bytes()

	fTime := time.Now()
	eTime := fTime.Sub(sTime)
	return "Image processed in " + fmt.Sprint(eTime), "jpeg.jpg", output, nil
}
