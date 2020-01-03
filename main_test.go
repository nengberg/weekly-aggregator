package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticate)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 302, rr.Code)
}
