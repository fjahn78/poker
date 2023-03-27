package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
  store := InMemoryPlayerStore{}
  server := PlayerServer{&store}
  player := "Pepper"

  server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player))
  server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player))
  server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player))

  response := httptest.NewRecorder()
  server.ServeHTTP(response, newGetScoreRequest(player))
  assertStatus(t, response.Code, http.StatusOK)

  assertResponseBody(t, response.Body.String(), "3")
}