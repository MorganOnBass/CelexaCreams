package celexacreams

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func FindNearestImageURL(m *discordgo.MessageCreate, c *discordgo.Channel, s *discordgo.Session) (string, error) {
	var url string
	var err error
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
