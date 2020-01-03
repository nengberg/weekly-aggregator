package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGetID(t *testing.T) {
	p := Playlist{
		URI: "abc:1:1",
	}

	result, err := p.GetID()

	assert.NoError(t, err)
	assert.Equal(t, "1", result)
}

func TestGetIDMalformed(t *testing.T) {
	var tests = []struct {
		URI string
	}{
		{"abc:"},
		{":abc"},
		{"abc"},
		{"a:bc"},
		{"a:bc:"},
	}

	for _, tt := range tests {
		t.Run(tt.URI, func(t *testing.T) {

			p := Playlist{
				URI: tt.URI,
			}

			_, err := p.GetID()

			assert.Error(t, err)

		})
	}
}

func TestGetPlaylists(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se"
	playlists := Playlists{
		Items: []Playlist{
			Playlist{Name: "first", URI: "playlist:123abc"},
			Playlist{Name: "second", URI: "playlist:123Ws"}},
	}
	jsonPlaylists, _ := json.Marshal(playlists)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonPlaylists))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}
	client.HTTP = cli

	result, err := client.GetPlaylists()

	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, playlists.Items[0].Name, result.Items[0].Name)
	assert.Equal(t, playlists.Items[0].URI, result.Items[0].URI)
	assert.Equal(t, playlists.Items[1].Name, result.Items[1].Name)
	assert.Equal(t, playlists.Items[1].URI, result.Items[1].URI)
}

func TestGetPlaylistsError(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}
	client.HTTP = cli

	result, err := client.GetPlaylists()

	assert.Error(t, err)
	assert.Empty(t, result.Items)
}

func TestAddTracksToPlaylist_TracksAlreadyInPlaylist_DoesntPostToPlaylists(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se/"
	playlistID := "123"
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/playlists/"+playlistID+"/tracks?offset=0" {
			tracksJSON, _ := json.Marshal(tracks)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tracksJSON))
		}

		assert.Equal(t, "/playlists/"+playlistID+"/tracks?offset=0", r.URL.String())
		assert.NotEqual(t, "/playlists/"+playlistID+"/tracks", r.URL.String())
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}
	client.HTTP = cli

	client.AddTracksToPlaylist(playlistID, tracks.Items)

}

func TestAddTracksToPlaylist_TracksNotInPlaylist_PostToPlaylists(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se/"
	playlistID := "123"
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/playlists/"+playlistID+"/tracks?offset=0" {
			tracksJSON, _ := json.Marshal(tracks)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tracksJSON))
		}

		if r.URL.String() == "/playlists/"+playlistID+"/tracks?offset=0" {
			//Test if post body is correct
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}
	client.HTTP = cli
	var tracks = TracksResponse{Items: []*Item{&Item{Track: Track{Name: "def", URI: "t:3:4"}}}}

	client.AddTracksToPlaylist(playlistID, tracks.Items)

}
