package handler

import (
	"bytes"
	"fmt"
	"git.qianqiusoft.com/library/graphics-go/graphics"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	"github.com/morganonbass/celexacreams/celexacreams"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"strconv"
	"sync"
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

// Help returns a brief help string
func (h *Spin) Help(short bool) string {
	if short {
		return "Crops a circle from the center of an image and returns a spinning gif. "
	} else {
		return fmt.Sprintf("Usage: `%vspin [1-10]`\n\nReturn: A spinning circle from the center of the supplied image. " +
			"You may attach an image or link to an image to the invoking post, or reply to a post containing an image. " +
			"If you do not, spin will process the most recent image in the channel.\n\n" +
			"The optional argument sets the spin speed.", celexacreams.Prefix)
	}
}


// Handle resizes to a sane value for the guild, crops to a circle, and returns a spinning gif
func (h *Spin) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	args, _ := celexacreams.ExtractCommand(m.Content)
	var speed int
	if len(args) > 1 {
		arg, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			arg = 5 // probably got '.spin <url>'
		}
		speed = int(arg)
	} else {
		speed = 5
	}
	if speed < 1 || speed > 10 {
		return "spin only accepts a speed between 1 and 10", "", make([]byte, 0), nil
	}
	delay := int(40 / float64(speed) - 2)
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
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	var resized *image.NRGBA
	// output gifs can get big, blowing past the upload limit of unboosted guilds
	if guild.PremiumTier < 2 {
		resized = imaging.Fit(i, 512, 512, imaging.Lanczos)
	} else {
		resized = imaging.Fit(i, 800, 800, imaging.Lanczos)
	}
	centre := []int{resized.Rect.Max.X / 2, resized.Rect.Max.Y / 2}
	palette := celexacreams.QuantizeImage(resized)
	images :=  make([]*image.Paletted, 36)
	delays := make([]int, 36)
	var wg sync.WaitGroup
	// spin this tasty record
	for f := 0; f <= 35; f++ {
		wg.Add(1)
		go drawFrame(&wg, resized, f, palette, centre, images, delays, delay)
	}
	wg.Wait()
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

func drawFrame(wg *sync.WaitGroup, resized *image.NRGBA, f int, palette color.Palette, centre []int, images []*image.Paletted, delays []int, delay int) {
	defer wg.Done()
	tmp := image.NewRGBA(resized.Bounds())
	graphics.Rotate(tmp, resized, &graphics.RotateOptions{(math.Pi / 18.0) * float64(f)})
	frame := image.NewPaletted(resized.Bounds(), palette)
	draw.DrawMask(frame, frame.Bounds(), tmp, image.Point{}, &celexacreams.Circle{
		P: image.Point{X: centre[0], Y: centre[1]},
		R: celexacreams.Min(resized.Rect.Max.X, resized.Rect.Max.Y) / 2,
	}, image.Point{}, draw.Over)
	images[f] = frame
	delays[f] = delay
}
