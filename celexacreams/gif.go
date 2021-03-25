package celexacreams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

type GiphySearchResponse struct {
	Data []GIFData `json:data`
}

type GIFData struct {
	URL string `json:url`
}

func GetGIF(search string) (string, error) {
	apiKey, ok := os.LookupEnv("GIPHY_API_KEY")
	if !ok {
		return "", fmt.Errorf("GIPHY_API_KEY not found")
	}

	offset := rand.Intn(50)
	resp, err := http.Get(
		fmt.Sprintf(
			"https://api.giphy.com/v1/gifs/search?api_key=%s&q=%s&limit=1&offset=%d&rating=G&lang=en",
			apiKey,
			url.QueryEscape(search),
			offset,
		),
	)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"HTTP status %d", resp.StatusCode,
		)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Info(bodyBytes)

	searchResults := &GiphySearchResponse{}
	err = json.Unmarshal(bodyBytes, searchResults)
	if err != nil {
		return "", err
	}

	return searchResults.Data[0].URL, nil
}
