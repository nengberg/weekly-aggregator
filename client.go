package main

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// Client is a wrapper working with Spotify API
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

//NewClient is a function to create a new client for communicating with Spotify API
func (a *Auth) NewClient(token *oauth2.Token) Client {
	http := a.Config.Client(context.Background(), token)
	return Client{
		BaseURL: "https://api.spotify.com/v1/",
		HTTP:    http,
	}
}
