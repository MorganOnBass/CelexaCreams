package celexacreams

import (
	"fmt"
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

func GetImageURLFromMessage(m *discordgo.Message) (string, error) {
	if len(m.Attachments) > 0 {
		url := m.Attachments[0].URL
		return url, nil
	}
	splitInput := strings.Split(m.ContentWithMentionsReplaced(), " ")
	for _, input := range splitInput {
		if IsURL(input) {
			resp, err := http.Head(input)
			if err != nil {
				continue
			}
			if resp.StatusCode != http.StatusOK {
				continue
			}
			if resp.Header.Get("content-type") == "image/gif" {
				continue
			}
			if strings.HasPrefix(resp.Header.Get("content-type"), "image") {
				return input, nil
			}

		}
	}
	return "", fmt.Errorf(
		"no Image Found",
	)
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
