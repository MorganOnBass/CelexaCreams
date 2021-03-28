package handler

import (
	"github.com/morganonbass/celexacreams/celexacreams"
	"gopkg.in/gographics/imagick.v3/imagick"

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

	//return "Not implemented yet, but probably logging to console for Morgan to debug", make([]byte, 0), nil
}
