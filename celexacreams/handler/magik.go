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
type Magik struct{}

// Handle meows back
func (h *Magik) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, []byte, error) {
	args, _ := celexacreams.ExtractCommand(m.Content)
	var sauce float64
	if len(args) > 1 {
		arg, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			sauce = float64(2.0) // probably got '.magik <URL>'
		}
		sauce = arg
	} else {
		sauce = float64(2.0)
	}
	if sauce == float64(1.0) {
		// this gets rounded to a uint for mw.LiquidRescaleImage() so 1.0 is no weaker than default
		sauce = float64(0.9)
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
		return "", make([]byte, 0), err
	}
	defer s.ChannelMessageDelete(c.ID, msg.ID)
	sTime := time.Now()
	URL, err := celexacreams.FindNearestImageURL(m, c, s)
	if err != nil {
		return "", make([]byte, 0), err
	}
	image, err := celexacreams.DownloadImage(URL)
	if err != nil {
		return "", make([]byte, 0), err
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err = mw.ReadImageBlob(image)
	if err != nil {
		return "", make([]byte, 0), err
	}
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	err = mw.LiquidRescaleImage(uint(float64(width)*0.5), uint(float64(height)*0.5), math.Round(sauce*0.5), 0)
	if err != nil {
		return "", make([]byte, 0), err
	}
	err = mw.LiquidRescaleImage(uint(float64(width)*0.75), uint(float64(height)*0.75), math.Round(sauce), 0)
	if err != nil {
		return "", make([]byte, 0), err
	}
	err = mw.SetImageFormat("PNG")
	if err != nil {
		return "", make([]byte, 0), err
	}
	output := mw.GetImageBlob()
	fTime := time.Now()
	eTime := fTime.Sub(sTime)
	return "Image processed in " + fmt.Sprint(eTime), output, nil
}
