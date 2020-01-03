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
func (p Playlist) GetID() (string, error) {
	uriSplitted := strings.Split(p.URI, ":")

	if len(uriSplitted) < 3 || uriSplitted[2] == "" {
		return "", fmt.Errorf("Invalid URI format %v", p.URI)
	}
	return uriSplitted[2], nil
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
	tracksAlreadyInList := getTracks(c, playlistID)
	tracksToAdd := findDeltas(tracksAlreadyInList, items)
	if len(tracksToAdd) != 0 {
		data, _ := createPostBody(tracksToAdd)
		res, err := c.HTTP.Post(playlistTracksURL, "application/json", bytes.NewReader(data))
		if err != nil {
			log.Fatalf("An error occurred when adding tracks %s\n", err.Error())
		}
		defer res.Body.Close()
	} else {
		log.Print("No tracks to add!")
	}
}

func getTracks(c *Client, playlistID string) []*Item {
	tracksResponse, _ := c.GetTracks(playlistID, 0)
	tracks := tracksResponse.Items
	total := tracksResponse.Total
	fetched := 100
	for {
		if fetched >= total {
			break
		}

		tr, _ := c.GetTracks(playlistID, fetched)
		tracks = append(tracks, tr.Items...)
		fetched += 100
	}
	log.Printf("Length: %d", len(tracks))
	return tracks
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
