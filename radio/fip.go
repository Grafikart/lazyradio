package radio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type fipResponse struct {
	Songs []struct {
		FirstLine  string `json:"firstLine"`
		SecondLine string `json:"secondLine"`
	} `json:"songs"`
}

func FetcherFipTracks(radio string) func() ([]TrackItem, error) {
	return func() ([]TrackItem, error) {
		client := http.Client{Timeout: 2 * time.Second}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://www.radiofrance.fr/fip/%v/api/songs", radio), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
		}
		var apiResponse fipResponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read response body : %v", err)
		}
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return nil, fmt.Errorf("cannot parse fip response : %v", err)
		}
		var tracks []TrackItem
		for _, song := range apiResponse.Songs {
			tracks = append(tracks, TrackItem{
				Name:   song.SecondLine,
				Artist: song.FirstLine,
			})
		}
		return tracks, nil
	}
}
