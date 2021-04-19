package handler

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/morganonbass/celexacreams/celexacreams"
	"gopkg.in/gographics/imagick.v3/imagick"

	"github.com/bwmarrin/discordgo"
)

// Magik responds to "magik"
type Magik struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *Magik) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Magik) DeleteInvocation() bool {
	return h.D
}

// Help returns a brief help string
func (h *Magik) Help(short bool) string {
	if short {
		return "Does magik on an image. Optional integer argument sets just how much magik to do"
	} else {
		return fmt.Sprintf("Usage: `%vmagik [integer]`\n\nYou may attach an image or link to an image to the invoking post," +
			" or reply to a post containing an image. If you do not, magik will process the most recent image in the channel.",
			celexacreams.Prefix)
	}
}


// Handle does some magik mangling. An optional argument can make it even more magikal
func (h *Magik) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	args, _ := celexacreams.ExtractCommand(m.Content)
	var sauce float64
	if len(args) > 1 {
		arg, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			arg = float64(2.0) // probably got '.magik <URL>'
		}
		sauce = arg
	} else {
		sauce = float64(2.0)
	}
	if sauce == float64(1.0) {
		// this gets rounded to a uint for mw.LiquidRescaleImage() so 1.0 is no weaker than default
		sauce = float64(0.9)
	}
	if sauce < 0 {
		// this segfaults imagemagick, lol
		return "Negative numbers make imagemagick segfault. Are you trying to kill me?", "", make([]byte, 0), nil
	}
	if sauce > 50000 {
		return "Nothing requires that much magik sauce.", "", make([]byte, 0), nil
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
	sTime := time.Now()
	URL, err := celexacreams.FindNearestImageURL(m, c, s)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	image, err := celexacreams.DownloadImage(URL)
	if err != nil {
		return "", "", make([]byte, 0), err
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err = mw.ReadImageBlob(image)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.AutoOrientImage()
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	err = mw.LiquidRescaleImage(uint(float64(width)*0.5), uint(float64(height)*0.5), math.Round(sauce*0.5), 0)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.LiquidRescaleImage(uint(float64(width)*0.75), uint(float64(height)*0.75), math.Round(sauce), 0)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.SetImageFormat("PNG")
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	output := mw.GetImageBlob()
	fTime := time.Now()
	eTime := fTime.Sub(sTime)
	return "Image processed in " + fmt.Sprint(eTime), "magik.png", output, nil
}
