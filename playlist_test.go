package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	client.HTTP = createNewHTTPServerWithClient(handler)

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
	client.HTTP = createNewHTTPServerWithClient(handler)

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
	tracksResponse := TracksResponse{Items: []*Item{&Item{Track: Track{Name: "abc", URI: "t:1:2"}}}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/playlists/"+playlistID+"/tracks?offset=0" {
			tracksJSON, _ := json.Marshal(tracksResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tracksJSON))
		}

		assert.Equal(t, "/playlists/"+playlistID+"/tracks?offset=0", r.URL.String())
		assert.Equal(t, "GET", r.Method)
		assert.NotEqual(t, "/playlists/"+playlistID+"/tracks", r.URL.String())
		assert.NotEqual(t, "POST", r.Method)

	}
	client.HTTP = createNewHTTPServerWithClient(handler)

	client.AddTracksToPlaylist(playlistID, tracks.Items)
}

func TestAddTracksToPlaylist_TracksNotInPlaylist_PostToPlaylists(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se/"
	playlistID := "123"
	tracks := TracksResponse{Items: []*Item{&Item{Track: Track{Name: "abc", URI: "t:1:2"}}}}
	tracksToAdd := []*Item{
		&Item{
			Track: Track{
				Name: "def",
				URI:  "t:3:4",
			},
		},
		&Item{
			Track: Track{
				Name: "ghi",
				URI:  "t:8:5",
			},
		},
	}
	tracksToAddResponse := TracksResponse{Items: tracksToAdd}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/playlists/"+playlistID+"/tracks?offset=0" && r.Method == "GET" {
			tracksJSON, _ := json.Marshal(tracks)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tracksJSON))
		}

		if r.URL.String() == "/playlists/"+playlistID+"/tracks" && r.Method == "POST" {
			responseData, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			request := make(map[string][]string)
			err = json.Unmarshal(responseData, &request)
			uris := request["uris"]
			assert.NoError(t, err)
			assert.Equal(t, 2, len(uris))
			assert.Equal(t, tracksToAddResponse.Items[0].Track.URI, uris[0])
			assert.Equal(t, tracksToAddResponse.Items[1].Track.URI, uris[1])
		}
	}
	client.HTTP = createNewHTTPServerWithClient(handler)

	client.AddTracksToPlaylist(playlistID, tracksToAddResponse.Items)
}

func TestAddTracksToPlaylist_AllTracksInList_MoreTracksThanDefaultPaging_DoesntPostToPlaylists(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se/"
	playlistID := "123"
	total := 101
	tracks := makeTracks(total)
	tracksResponse := TracksResponse{Items: tracks, Total: total}
	var urlsCalled []string
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tracksJSON, _ := json.Marshal(tracksResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(tracksJSON))
		}
		urlsCalled = append(urlsCalled, r.URL.String())

		assert.NotEqual(t, "/playlists/"+playlistID+"/tracks", r.URL.String())
		assert.NotEqual(t, "POST", r.Method)
	}
	client.HTTP = createNewHTTPServerWithClient(handler)

	client.AddTracksToPlaylist(playlistID, tracksResponse.Items)

	assert.Equal(t, 2, len(urlsCalled))
	assert.Equal(t, "/playlists/"+playlistID+"/tracks?offset=0", urlsCalled[0])
	assert.Equal(t, "/playlists/"+playlistID+"/tracks?offset=100", urlsCalled[1])
}

func makeTracks(total int) []*Item {
	items := make([]*Item, total)
	for i := 0; i < total; i++ {
		track := &Item{Track: Track{Name: fmt.Sprintf("track-%d", i), URI: fmt.Sprintf("uri-%d", i)}}
		items[i] = track
	}
	return items
}

func createNewHTTPServerWithClient(handler func(w http.ResponseWriter, r *http.Request)) *http.Client {
	server := httptest.NewServer(http.HandlerFunc(handler))
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}
}
