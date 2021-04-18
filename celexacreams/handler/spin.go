package handler

import (
	"bytes"
	"fmt"
	"git.qianqiusoft.com/library/graphics-go/graphics"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	"github.com/morganonbass/celexacreams/celexacreams"
	"image"
	"image/draw"
	"image/gif"
	"math"
	"time"
)

// Spin responds to "spin"
type Spin struct {
	R, D bool
}

// Should I reply to the invoking message?
func (h *Spin) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Spin) DeleteInvocation() bool {
	return h.D
}

// Handle resizes to 256x256, crops to a circle, and creates a spinning gif
func (h *Spin) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	sTime := time.Now()
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
	URL, err := celexacreams.FindNearestImageURL(m, c, s)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	img, err := celexacreams.DownloadImage(URL)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	i, err := imaging.Decode(bytes.NewReader(img), imaging.AutoOrientation(true))
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	// cut input down to a sane size so this doesn't take all day and make a huge gif
	resized := imaging.Fit(i, 800, 800, imaging.Lanczos)
	centre := []int{resized.Rect.Max.X / 2, resized.Rect.Max.Y / 2}
	var images []*image.Paletted
	var delays []int
	for f := 0; f <= 35; f++ {
		tmp := image.NewRGBA(resized.Bounds())
		graphics.Rotate(tmp, resized, &graphics.RotateOptions{(math.Pi / 18.0) * float64(f)})
		frame := image.NewPaletted(resized.Bounds(), celexacreams.Palette)
		draw.DrawMask(frame, frame.Bounds(), tmp, image.Point{}, &celexacreams.Circle{
			P: image.Point{X: centre[0], Y: centre[1]},
			R: celexacreams.Min(resized.Rect.Max.X, resized.Rect.Max.Y) / 2,
		}, image.Point{}, draw.Over)
		images = append(images, frame)
		delays = append(delays, 0)
	}
	buf := new(bytes.Buffer)
	err = gif.EncodeAll(buf, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	output := buf.Bytes()

	fTime := time.Now()
	eTime := fTime.Sub(sTime)
	return "Image processed in " + fmt.Sprint(eTime), "test.gif", output, nil

}
