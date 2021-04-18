package celexacreams

import (
	"fmt"
	"github.com/soniakeys/quant/median"
	"image"
	"image/color"
	"image/draw"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func FindNearestImageURL(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, error) {
	var url string
	var err error
	if m.MessageReference != nil {
		// We seem to have been invoked with a reply/crosspost
		ref, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		if err == nil {
			url, err = GetImageURLFromMessage(ref)
			if err == nil {
				return url, nil
			}
		}
		// Something went wrong retrieving an image URL from the message reference, carry on as normal...
	}
	url, err = GetImageURLFromMessage(m.Message)
	if err != nil {
		history, err := s.ChannelMessages(c.ID, 100, m.ID, "", "")
		if err != nil {
			return "", fmt.Errorf("error retrieving message history")
		}
		for _, msg := range history {
			url, err = GetImageURLFromMessage(msg)
			if err != nil {
				continue
			}
			return url, nil
		}
	}
	if err != nil {
		return "", err
	}
	return url, nil
}

func IsImage(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP Status %d", resp.StatusCode)
	}
	if resp.Header.Get("content-type") == "image/gif" {
		return false, nil
	}
	if strings.HasPrefix(resp.Header.Get("content-type"), "image") {
		return true, nil
	}
	return false, nil
}

func GetImageURLFromMessage(m *discordgo.Message) (string, error) {
	if len(m.Attachments) > 0 {
		url := m.Attachments[0].URL
		isImage, err := IsImage(url)
		if err != nil || !isImage {
			return "", &NotAnImageError{url}
		}
		return url, nil
	}
	splitInput := strings.Split(m.ContentWithMentionsReplaced(), " ")
	for _, url := range splitInput {
		if IsURL(url) {
			isImage, err := IsImage(url)
			if err != nil || !isImage {
				return "", &NotAnImageError{url}
			}
			return url, nil
		}
	}
	return "", &ImageNotFoundError{m.ID}
}

func DownloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf(
			"HTTP status %d", resp.StatusCode,
		)
	}
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return image, nil
}

type Pixel struct {
	Point image.Point
	Color color.Color
}

// Decode image.Image's pixel data into []*Pixel
func DecodePixelsFromImage(img image.Image, offsetX, offsetY int) []*Pixel {
	pixels := []*Pixel{}
	for y := 0; y <= img.Bounds().Max.Y; y++ {
		for x := 0; x <= img.Bounds().Max.X; x++ {
			p := &Pixel{
				Point: image.Point{X: x + offsetX, Y: y + offsetY},
				Color: img.At(x, y),
			}
			pixels = append(pixels, p)
		}
	}
	return pixels
}

type Circle struct {
	P image.Point
	R int
}

func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.P.X-c.R, c.P.Y-c.R, c.P.X+c.R, c.P.Y+c.R)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.P.X)+0.5, float64(y-c.P.Y)+0.5, float64(c.R)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{A: 255}
	}
	return color.Alpha{}
}

func DrawTransparentLayer(res []int) draw.Image {
	backGroundColor := image.Transparent
	backgroundWidth := res[0]
	backgroundHeight := res[1]
	background := image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))

	draw.Draw(background, background.Bounds(), backGroundColor, image.Point{}, draw.Src)
	return background
}

func QuantizeImage(i image.Image) color.Palette {
	p := make([]color.Color, 0, 255)
	var palette = color.Palette{
		image.Transparent,
	}
	var q draw.Quantizer = median.Quantizer(255)
	p = q.Quantize(p, i)
	for _, c := range p {
		palette = append(palette, c)
	}
	return palette
}
