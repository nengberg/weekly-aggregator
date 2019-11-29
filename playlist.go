package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Playlists holds a list of playlists
type Playlists struct {
	Items []Playlist `json:"items"`
}

//Playlist information
type Playlist struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

//GetID returns an ID for a Playlist
func (p Playlist) GetID() string {
	uriSplitted := strings.Split(p.URI, ":")
	return uriSplitted[2]
}

//GetPlaylists returns your playlists
func (c *Client) GetPlaylists() (Playlists, error) {
	playlistsURL := c.BaseURL + "me/playlists"
	response, err := c.HTTP.Get(playlistsURL)
	if err != nil {
		return Playlists{}, err
	}
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	var playlists Playlists
	if err := json.Unmarshal(bodyBytes, &playlists); err != nil {
		return Playlists{}, err
	}
	return playlists, nil
}

//AddTracksToPlaylist adds the given track list to the playlist matching playlistID
func (c *Client) AddTracksToPlaylist(playlistID string, items []*Item) {
	playlistTracksURL := c.BaseURL + "playlists/" + playlistID + "/tracks"
	log.Printf("Adding tracks to playlist with ID %s\n", playlistTracksURL)
	tracksAlreadyInList, _ := c.GetTracks(playlistID)
	tracksToAdd := findDeltas(tracksAlreadyInList, items)
	data, _ := createPostBody(tracksToAdd)
	res, err := c.HTTP.Post(playlistTracksURL, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("An error occurred when adding tracks %s\n", err.Error())
	}
	defer res.Body.Close()
}

func findDeltas(tracksAlreadyInList []*Item, tracksToAdd []*Item) []*Item {

	var unique []*Item

	for _, item := range tracksToAdd {
		skip := false
		for _, u := range tracksAlreadyInList {
			if item.Track.URI == u.Track.URI {
				fmt.Printf("Skipping track '%s' as it is already in playlist.\n", item.Track.Name)
				skip = true
				break
			}
		}
		if !skip {
			unique = append(unique, item)
		}
	}

	return unique
}

func createPostBody(items []*Item) ([]byte, error) {
	uris := make([]string, len(items))
	for i, item := range items {
		fmt.Printf("%+v\n", item.Track.Name)
		uris[i] = item.Track.URI
	}
	m := make(map[string]interface{})
	m["uris"] = uris
	return json.Marshal(m)
}
