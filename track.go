package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// TracksResponse is a container for JSON response
type TracksResponse struct {
	Items []*Item `json:"items"`
	Total int     `json:"total"`
}

//Item is a container holding tracks from a Spotify playlist
type Item struct {
	Track Track `json:"track"`
}

//Track is the object for a playlist track
type Track struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

//GetTracks fetches tracks from a given playlistID
func (c *Client) GetTracks(playlistID string, offset int) (TracksResponse, error) {
	playlistTracksURL := fmt.Sprintf("%splaylists/%s/tracks?offset=%d", c.BaseURL, playlistID, offset)
	response, err := c.HTTP.Get(playlistTracksURL)
	if err != nil {
		log.Fatalf("An error occurred when getting tracks from url %s %s", playlistTracksURL, err.Error())
		return TracksResponse{}, nil
	}
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	var resp TracksResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return TracksResponse{}, err
	}
	return resp, nil
}
