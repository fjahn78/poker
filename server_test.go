package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetPlayerScore("Pepper")
		response := httptest.NewRecorder()

		want := "20"

		PlayerServer(response, request)

		assertResponseBody(t, response, want)
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetPlayerScore("Floyd")
		response := httptest.NewRecorder()

		want := "10"

		PlayerServer(response, request)

		assertResponseBody(t, response, want)
	})
}

func assertResponseBody(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Body.String()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func newGetPlayerScore(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}
