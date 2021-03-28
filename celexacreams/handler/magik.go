package handler

import (
	"github.com/morganonbass/celexacreams/celexacreams"
	"gopkg.in/gographics/imagick.v3/imagick"

	"github.com/bwmarrin/discordgo"
)

// Magik responds to "magik"
type Magik struct{}

// Handle meows back
func (h *Magik) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, []byte, error) {
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

	err = mw.LiquidRescaleImage(uint(float64(width)*0.5), uint(float64(height)*0.5), 1, 0)
	if err != nil {
		return "", make([]byte, 0), err
	}
	err = mw.LiquidRescaleImage(uint(float64(width)*0.75), uint(float64(height)*0.75), 2, 0)
	if err != nil {
		return "", make([]byte, 0), err
	}
	err = mw.SetImageFormat("PNG")
	if err != nil {
		return "", make([]byte, 0), err
	}
	output := mw.GetImageBlob()
	return "", output, nil
}
