package handler

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// DeepFry responds to "deepfry"
type DeepFry struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *DeepFry) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *DeepFry) DeleteInvocation() bool {
	return h.D
}

// Handle returns a deep fried image
func (h *DeepFry) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
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


	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(img)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.AutoOrientImage()
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.SharpenImage(0.0, 10.0)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.ModulateImage(100.0, 650.0, 100.0)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw.BrightnessContrastImage(10.0, 47.5)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	args := []float64{7.379,22.860,1.870,2.611,-0.000}
	err = mw.FunctionImage(imagick.FUNCTION_POLYNOMIAL, args)
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
	return "Image processed in " + fmt.Sprint(eTime), "deepfried.png", output, nil
}
