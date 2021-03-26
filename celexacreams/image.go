package celexacreams

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func FindNearestImage(m *discordgo.MessageCreate) (string, error) {
	if len(m.Attachments) > 0 {
		url := m.Attachments[0].URL
		return url, nil
	}
	splitInput := strings.Split(m.ContentWithMentionsReplaced(), " ")
	for _, input := range splitInput {
		if IsURL(input) {
			resp, err := http.Head(input)
			if err != nil {
				panic(err)
			}
			if resp.StatusCode != http.StatusOK {
				return "", fmt.Errorf(
					"HTTP status %d", resp.StatusCode,
				)
			}
			if strings.HasPrefix(resp.Header.Get("content-type"), "image") {
				return input, nil
			}

		}
	}
	return "", fmt.Errorf(
		"No Image Found.",
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
