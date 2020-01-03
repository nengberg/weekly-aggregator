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

func TestGetTracks(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)
	client.BaseURL = "http://abc.se"
	jsonTracks, _ := json.Marshal(tracks)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonTracks))
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

	result, err := client.GetTracks("playlistID", 0)

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
}

func TestGetTracksError(t *testing.T) {
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

	result, err := client.GetTracks("playlistID", 0)

	assert.Error(t, err)
	assert.Empty(t, result.Items)
}

var tracks = TracksResponse{Items: []*Item{&Item{Track: Track{Name: "abc", URI: "t:1:2"}}}}
