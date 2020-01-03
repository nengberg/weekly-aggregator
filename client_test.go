package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	auth := NewAuth()
	token := &oauth2.Token{}
	client := auth.NewClient(token)

	assert.Equal(t, "https://api.spotify.com/v1/", client.BaseURL)
	assert.NotNil(t, client.HTTP)
}
