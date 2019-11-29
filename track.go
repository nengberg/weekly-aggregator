package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// TracksResponse is a container for JSON response
type TracksResponse struct {
	Items []*Item `json:"items"`
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
func (c *Client) GetTracks(playlistID string) ([]*Item, error) {
	playlistTracksURL := c.BaseURL + "playlists/" + playlistID + "/tracks"
	response, err := c.HTTP.Get(playlistTracksURL)
	if err != nil {
		log.Fatalf("An error occurred when getting tracks from url %s %s", playlistTracksURL, err.Error())
		return []*Item{}, nil
	}
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	var resp TracksResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return []*Item{}, err
	}
	return resp.Items, nil
}
