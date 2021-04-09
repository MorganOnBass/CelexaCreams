package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/morganonbass/celexacreams/celexacreams"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// Ahha responds to "ahha"
type Ahha struct{
	R, D bool
}

// Should I reply to the invoking message?
func (h *Ahha) Reply() bool {
	return h.R
}

// Should I delete the invoking message?
func (h *Ahha) DeleteInvocation() bool {
	return h.D
}

// Handle returns an image mirrored about the Y axis
func (h *Ahha) Handle(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, string, []byte, error) {
	sTime := time.Now()
	URL, err := celexacreams.FindNearestImageURL(m, c, s)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	img, err := celexacreams.DownloadImage(URL)
	if err != nil {
		return "", "", make([]byte, 0), err
	}

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
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	mw2 := mw.Clone()
	defer mw2.Destroy()
	err = mw2.CropImage(width/2, height, int(width/2), 0)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	mw1 := mw2.Clone()
	defer mw1.Destroy()
	err = mw1.RotateImage(imagick.NewPixelWand(), float64(180))
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw1.FlipImage()
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw1.SetImageFormat("PNG")
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	err = mw2.SetImageFormat("PNG")
	if err != nil {
		return "", "", make([]byte, 0), err
	}

	h1blob := mw1.GetImageBlob()
	h2blob := mw2.GetImageBlob()
	h1, _, err := image.Decode(bytes.NewReader(h1blob))
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	h2, _, err := image.Decode(bytes.NewReader(h2blob))
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	pixels1 := celexacreams.DecodePixelsFromImage(h1, 0, 0)
	pixels2 := celexacreams.DecodePixelsFromImage(h2, h1.Bounds().Max.X, 0)
	pixelSum := append(pixels1, pixels2...)
	newRect := image.Rectangle{
		Min: h1.Bounds().Min,
		Max: image.Point{
			X: h2.Bounds().Max.X + h1.Bounds().Max.X,
			Y: h2.Bounds().Max.Y,
		},
	}
	finImage := image.NewRGBA(newRect)
	for _, px := range pixelSum {
		finImage.Set(
			px.Point.X,
			px.Point.Y,
			px.Color,
		)
	}
	draw.Draw(finImage, finImage.Bounds(), finImage, image.Point{0, 0}, draw.Src)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, finImage)
	if err != nil {
		return "", "", make([]byte, 0), err
	}
	output := buf.Bytes()

	fTime := time.Now()
	eTime := fTime.Sub(sTime)
	return "Image processed in " + fmt.Sprint(eTime), "ahha.png", output, nil
}
